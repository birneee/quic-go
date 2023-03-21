package quic

import (
	"context"
	"errors"
	"fmt"
	"github.com/lucas-clemente/quic-go/handover"
	"github.com/lucas-clemente/quic-go/internal/xse"
	"net"
	"sync"

	"github.com/lucas-clemente/quic-go/internal/flowcontrol"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type streamError struct {
	message string
	nums    []protocol.StreamNum
}

func (e streamError) Error() string {
	return e.message
}

func convertStreamError(err error, stype protocol.StreamType, pers protocol.Perspective) error {
	strError, ok := err.(streamError)
	if !ok {
		return err
	}
	ids := make([]interface{}, len(strError.nums))
	for i, num := range strError.nums {
		ids[i] = num.StreamID(stype, pers)
	}
	return fmt.Errorf(strError.Error(), ids...)
}

type streamOpenErr struct{ error }

var _ net.Error = &streamOpenErr{}

func (e streamOpenErr) Temporary() bool { return e.error == errTooManyOpenStreams }
func (streamOpenErr) Timeout() bool     { return false }

// errTooManyOpenStreams is used internally by the outgoing streams maps.
var errTooManyOpenStreams = errors.New("too many open streams")

type streamsMap struct {
	perspective protocol.Perspective
	version     protocol.VersionNumber

	maxIncomingBidiStreams uint64
	maxIncomingUniStreams  uint64

	sender            streamSender
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController

	mutex               sync.Mutex
	outgoingBidiStreams *outgoingStreamsMap[streamI]
	outgoingUniStreams  *outgoingStreamsMap[sendStreamI]
	incomingBidiStreams *incomingStreamsMap[streamI]
	incomingUniStreams  *incomingStreamsMap[receiveStreamI]
	reset               bool
	// if nil, XSE-QUIC extension is not used
	xseCryptoSetup       xse.CryptoSetup
	streamFramesInFlight func(streamID StreamID, encLevel protocol.EncryptionLevel) []*wire.StreamFrame
}

var _ streamManager = &streamsMap{}

func newStreamsMap(
	sender streamSender,
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController,
	maxIncomingBidiStreams uint64,
	maxIncomingUniStreams uint64,
	perspective protocol.Perspective,
	version protocol.VersionNumber,
	// required for stream state serialization
	streamFramesInFlight func(streamID StreamID, encLevel protocol.EncryptionLevel) []*wire.StreamFrame,
) streamManager {
	m := &streamsMap{
		perspective:            perspective,
		newFlowController:      newFlowController,
		maxIncomingBidiStreams: maxIncomingBidiStreams,
		maxIncomingUniStreams:  maxIncomingUniStreams,
		sender:                 sender,
		version:                version,
		streamFramesInFlight:   streamFramesInFlight,
	}
	m.initMaps()
	return m
}

func (m *streamsMap) initMaps() {
	m.outgoingBidiStreams = newOutgoingStreamsMap(
		protocol.StreamTypeBidi,
		func(num protocol.StreamNum) streamI {
			id := num.StreamID(protocol.StreamTypeBidi, m.perspective)
			stream := newStream(
				id,
				m.sender,
				m.newFlowController(id),
				m.version,
				func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
					return m.streamFramesInFlight(id, encLevel)
				},
			)
			if m.xseCryptoSetup != nil {
				return xseStreamI{m.xseCryptoSetup.NewStream(stream)}
			} else {
				return stream
			}
		},
		m.sender.queueControlFrame,
	)
	m.incomingBidiStreams = newIncomingStreamsMap(
		protocol.StreamTypeBidi,
		func(num protocol.StreamNum) streamI {
			id := num.StreamID(protocol.StreamTypeBidi, m.perspective.Opposite())
			stream := newStream(
				id,
				m.sender,
				m.newFlowController(id),
				m.version,
				func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
					return m.streamFramesInFlight(id, encLevel)
				},
			)
			if m.xseCryptoSetup != nil {
				return xseStreamI{m.xseCryptoSetup.NewStream(stream)}
			} else {
				return stream
			}
		},
		m.maxIncomingBidiStreams,
		m.sender.queueControlFrame,
	)
	m.outgoingUniStreams = newOutgoingStreamsMap(
		protocol.StreamTypeUni,
		func(num protocol.StreamNum) sendStreamI {
			id := num.StreamID(protocol.StreamTypeUni, m.perspective)
			stream := newSendStream(
				id,
				m.sender,
				m.newFlowController(id),
				m.version,
				func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
					return m.streamFramesInFlight(id, encLevel)
				},
			)
			if m.xseCryptoSetup != nil {
				return xseSendStreamI{m.xseCryptoSetup.NewSendStream(stream)}
			} else {
				return stream
			}
		},
		m.sender.queueControlFrame,
	)
	m.incomingUniStreams = newIncomingStreamsMap(
		protocol.StreamTypeUni,
		func(num protocol.StreamNum) receiveStreamI {
			id := num.StreamID(protocol.StreamTypeUni, m.perspective.Opposite())
			if m.xseCryptoSetup != nil {
				return xseReceiveStreamI{m.xseCryptoSetup.NewReceiveStream(newReceiveStream(id, m.sender, m.newFlowController(id), m.version))}
			} else {
				return newReceiveStream(id, m.sender, m.newFlowController(id), m.version)
			}
		},
		m.maxIncomingUniStreams,
		m.sender.queueControlFrame,
	)
}

func (m *streamsMap) OpenStream() (Stream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.outgoingBidiStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.OpenStream()
	return str, convertStreamError(err, protocol.StreamTypeBidi, m.perspective)
}

func (m *streamsMap) OpenStreamSync(ctx context.Context) (Stream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.outgoingBidiStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.OpenStreamSync(ctx)
	return str, convertStreamError(err, protocol.StreamTypeBidi, m.perspective)
}

func (m *streamsMap) OpenUniStream() (SendStream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.outgoingUniStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.OpenStream()
	return str, convertStreamError(err, protocol.StreamTypeBidi, m.perspective)
}

func (m *streamsMap) OpenUniStreamSync(ctx context.Context) (SendStream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.outgoingUniStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.OpenStreamSync(ctx)
	return str, convertStreamError(err, protocol.StreamTypeUni, m.perspective)
}

func (m *streamsMap) AcceptStream(ctx context.Context) (Stream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.incomingBidiStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.AcceptStream(ctx)
	return str, convertStreamError(err, protocol.StreamTypeBidi, m.perspective.Opposite())
}

func (m *streamsMap) AcceptUniStream(ctx context.Context) (ReceiveStream, error) {
	m.mutex.Lock()
	reset := m.reset
	mm := m.incomingUniStreams
	m.mutex.Unlock()
	if reset {
		return nil, Err0RTTRejected
	}
	str, err := mm.AcceptStream(ctx)
	return str, convertStreamError(err, protocol.StreamTypeUni, m.perspective.Opposite())
}

func (m *streamsMap) DeleteStream(id protocol.StreamID) error {
	num := id.StreamNum()
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			return convertStreamError(m.outgoingUniStreams.DeleteStream(num), protocol.StreamTypeUni, m.perspective)
		}
		return convertStreamError(m.incomingUniStreams.DeleteStream(num), protocol.StreamTypeUni, m.perspective.Opposite())
	case protocol.StreamTypeBidi:
		if id.InitiatedBy() == m.perspective {
			return convertStreamError(m.outgoingBidiStreams.DeleteStream(num), protocol.StreamTypeBidi, m.perspective)
		}
		return convertStreamError(m.incomingBidiStreams.DeleteStream(num), protocol.StreamTypeBidi, m.perspective.Opposite())
	}
	panic("")
}

func (m *streamsMap) GetOrOpenReceiveStream(id protocol.StreamID) (receiveStreamI, error) {
	str, err := m.getOrOpenReceiveStream(id)
	if err != nil {
		return nil, &qerr.TransportError{
			ErrorCode:    qerr.StreamStateError,
			ErrorMessage: err.Error(),
		}
	}
	return str, nil
}

func (m *streamsMap) getOrOpenReceiveStream(id protocol.StreamID) (receiveStreamI, error) {
	num := id.StreamNum()
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			// an outgoing unidirectional stream is a send stream, not a receive stream
			return nil, fmt.Errorf("peer attempted to open receive stream %d", id)
		}
		str, err := m.incomingUniStreams.GetOrOpenStream(num)
		return str, convertStreamError(err, protocol.StreamTypeUni, m.perspective)
	case protocol.StreamTypeBidi:
		var str receiveStreamI
		var err error
		if id.InitiatedBy() == m.perspective {
			str, err = m.outgoingBidiStreams.GetStream(num)
		} else {
			str, err = m.incomingBidiStreams.GetOrOpenStream(num)
		}
		return str, convertStreamError(err, protocol.StreamTypeBidi, id.InitiatedBy())
	}
	panic("")
}

func (m *streamsMap) GetOrOpenSendStream(id protocol.StreamID) (sendStreamI, error) {
	str, err := m.getOrOpenSendStream(id)
	if err != nil {
		return nil, &qerr.TransportError{
			ErrorCode:    qerr.StreamStateError,
			ErrorMessage: err.Error(),
		}
	}
	return str, nil
}

func (m *streamsMap) getOrOpenSendStream(id protocol.StreamID) (sendStreamI, error) {
	num := id.StreamNum()
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			str, err := m.outgoingUniStreams.GetStream(num)
			return str, convertStreamError(err, protocol.StreamTypeUni, m.perspective)
		}
		// an incoming unidirectional stream is a receive stream, not a send stream
		return nil, fmt.Errorf("peer attempted to open send stream %d", id)
	case protocol.StreamTypeBidi:
		var str sendStreamI
		var err error
		if id.InitiatedBy() == m.perspective {
			str, err = m.outgoingBidiStreams.GetStream(num)
		} else {
			str, err = m.incomingBidiStreams.GetOrOpenStream(num)
		}
		return str, convertStreamError(err, protocol.StreamTypeBidi, id.InitiatedBy())
	}
	panic("")
}

func (m *streamsMap) HandleMaxStreamsFrame(f *wire.MaxStreamsFrame) {
	switch f.Type {
	case protocol.StreamTypeUni:
		m.outgoingUniStreams.SetMaxStream(f.MaxStreamNum)
	case protocol.StreamTypeBidi:
		m.outgoingBidiStreams.SetMaxStream(f.MaxStreamNum)
	}
}

func (m *streamsMap) UpdateLimits(p *wire.TransportParameters) {
	m.outgoingBidiStreams.UpdateSendWindow(p.InitialMaxStreamDataBidiRemote)
	m.outgoingBidiStreams.SetMaxStream(p.MaxBidiStreamNum)
	m.outgoingUniStreams.UpdateSendWindow(p.InitialMaxStreamDataUni)
	m.outgoingUniStreams.SetMaxStream(p.MaxUniStreamNum)
}

func (m *streamsMap) CloseWithError(err error) {
	m.outgoingBidiStreams.CloseWithError(err)
	m.outgoingUniStreams.CloseWithError(err)
	m.incomingBidiStreams.CloseWithError(err)
	m.incomingUniStreams.CloseWithError(err)
}

// ResetFor0RTT resets is used when 0-RTT is rejected. In that case, the streams maps are
// 1. closed with an Err0RTTRejected, making calls to Open{Uni}Stream{Sync} / Accept{Uni}Stream return that error.
// 2. reset to their initial state, such that we can immediately process new incoming stream data.
// Afterwards, calls to Open{Uni}Stream{Sync} / Accept{Uni}Stream will continue to return the error,
// until UseResetMaps() has been called.
func (m *streamsMap) ResetFor0RTT() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.reset = true
	m.CloseWithError(Err0RTTRejected)
	m.initMaps()
}

func (m *streamsMap) UseResetMaps() {
	m.mutex.Lock()
	m.reset = false
	m.mutex.Unlock()
}

func (m *streamsMap) SetXseCryptoSetup(xseCryptoSetup xse.CryptoSetup) {
	m.xseCryptoSetup = xseCryptoSetup
}

func (m *streamsMap) RestoreBidiStream(state *handover.BidiStreamState) (Stream, error) {
	var stream Stream
	var err error
	if state.ID.InitiatedBy() == m.perspective {
		stream, err = m.outgoingBidiStreams.RestoreStream(state.ID.StreamNum(), state, m.perspective)
	} else {
		stream, err = m.incomingBidiStreams.RestoreStream(state.ID.StreamNum(), state, m.perspective)
	}
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func streamIToBidiStreamState(s streamI, perspective Perspective, config *ConnectionStateStoreConf) handover.BidiStreamState {
	ss := handover.BidiStreamState{
		ID: s.StreamID(),
	}
	s.storeReceiveState(&ss, perspective, config)
	s.storeSendState(&ss, perspective, config)
	return ss
}

func (m *streamsMap) BidiStreamStates(config *ConnectionStateStoreConf) map[StreamID]handover.BidiStreamState {
	states := make(map[protocol.StreamID]handover.BidiStreamState)
	for _, stream := range m.outgoingBidiStreams.streams {
		states[stream.StreamID()] = streamIToBidiStreamState(stream, m.perspective, config)
	}
	for _, entry := range m.incomingBidiStreams.streams {
		stream := entry.stream
		states[stream.StreamID()] = streamIToBidiStreamState(stream, m.perspective, config)
	}
	return states
}

func (m *streamsMap) OpenedBidiStream(id StreamID) (Stream, error) {
	if id.InitiatedBy() == m.perspective {
		return m.outgoingBidiStreams.GetStream(id.StreamNum())
	} else {
		return m.incomingBidiStreams.GetStream(id.StreamNum())
	}
}

func (m *streamsMap) UniStreamStates() map[protocol.StreamID]handover.UniStreamState {
	ss := make(map[protocol.StreamID]handover.UniStreamState)
	for range m.outgoingUniStreams.streams {
		//TODO implement me
		panic("implement me")
	}
	for range m.incomingUniStreams.streams {
		//TODO implement me
		panic("implement me")
	}
	return ss
}
