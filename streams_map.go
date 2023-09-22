package quic

import (
	"context"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/logging"
	"net"
	"sync"

	"github.com/quic-go/quic-go/internal/flowcontrol"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/wire"
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

	maxIncomingBidiStreams uint64
	maxIncomingUniStreams  uint64

	sender            streamSender
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController

	mutex                sync.Mutex
	outgoingBidiStreams  *outgoingStreamsMap[streamI]
	outgoingUniStreams   *outgoingStreamsMap[sendStreamI]
	incomingBidiStreams  *incomingStreamsMap[streamI]
	incomingUniStreams   *incomingStreamsMap[receiveStreamI]
	reset                bool
	streamFramesInFlight func(streamID StreamID, encLevel protocol.EncryptionLevel) []*wire.StreamFrame
}

var _ streamManager = &streamsMap{}

func newStreamsMap(
	sender streamSender,
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController,
	maxIncomingBidiStreams uint64,
	maxIncomingUniStreams uint64,
	perspective protocol.Perspective,
	// required for stream state serialization
	streamFramesInFlight func(streamID StreamID, encLevel protocol.EncryptionLevel) []*wire.StreamFrame,
) streamManager {
	m := &streamsMap{
		perspective:            perspective,
		newFlowController:      newFlowController,
		maxIncomingBidiStreams: maxIncomingBidiStreams,
		maxIncomingUniStreams:  maxIncomingUniStreams,
		sender:                 sender,
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
			return newStream(
				id,
				m.sender,
				m.newFlowController(id),
				func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
					return m.streamFramesInFlight(id, encLevel)
				},
			)
		},
		m.sender.queueControlFrame,
	)
	m.incomingBidiStreams = newIncomingStreamsMap(
		protocol.StreamTypeBidi,
		func(num protocol.StreamNum) streamI {
			id := num.StreamID(protocol.StreamTypeBidi, m.perspective.Opposite())
			return newStream(id, m.sender, m.newFlowController(id), func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
				return m.streamFramesInFlight(id, encLevel)
			})
		},
		m.maxIncomingBidiStreams,
		m.sender.queueControlFrame,
	)
	m.outgoingUniStreams = newOutgoingStreamsMap(
		protocol.StreamTypeUni,
		func(num protocol.StreamNum) sendStreamI {
			id := num.StreamID(protocol.StreamTypeUni, m.perspective)
			return newSendStream(id, m.sender, m.newFlowController(id), func(encLevel protocol.EncryptionLevel) []*wire.StreamFrame {
				return m.streamFramesInFlight(id, encLevel)
			})
		},
		m.sender.queueControlFrame,
	)
	m.incomingUniStreams = newIncomingStreamsMap(
		protocol.StreamTypeUni,
		func(num protocol.StreamNum) receiveStreamI {
			id := num.StreamID(protocol.StreamTypeUni, m.perspective.Opposite())
			return newReceiveStream(id, m.sender, m.newFlowController(id))
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

func (m *streamsMap) RestoreBidiStream(streamID StreamID, state *handover.BidiStreamState) (Stream, error) {
	var stream Stream
	var err error
	if streamID.InitiatedBy() == m.perspective {
		stream, err = RestoreOutgoingBidiStream(m.outgoingBidiStreams, streamID.StreamNum(), state, m.perspective)
	} else {
		stream, err = RestoreIncomingBidiStream(m.incomingBidiStreams, streamID.StreamNum(), state, m.perspective)
	}
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (m *streamsMap) RestoreSendStream(streamID StreamID, state *handover.UniStreamState) (SendStream, error) {
	return RestoreOutgoingUniStream(m.outgoingUniStreams, streamID.StreamNum(), state, m.perspective)
}

func (m *streamsMap) RestoreReceiveStream(streamID StreamID, state *handover.UniStreamState) (ReceiveStream, error) {
	return RestoreIncomingUniStream(m.incomingUniStreams, streamID.StreamNum(), state, m.perspective)
}

func streamIToBidiStreamState(s streamI, perspective logging.Perspective, config *ConnectionStateStoreConf) handover.BidiStreamState {
	ss := handover.BidiStreamState{}
	s.storeReceiveState(&ss, perspective, config)
	s.storeSendState(&ss, perspective, config)
	return ss
}

func (m *streamsMap) BidiStreamStates(config *ConnectionStateStoreConf) map[StreamID]*handover.BidiStreamState {
	states := make(map[protocol.StreamID]*handover.BidiStreamState)
	for _, stream := range m.outgoingBidiStreams.streams {
		state := streamIToBidiStreamState(stream, m.perspective, config)
		states[stream.StreamID()] = &state
	}
	for _, entry := range m.incomingBidiStreams.streams {
		stream := entry.stream
		state := streamIToBidiStreamState(stream, m.perspective, config)
		states[stream.StreamID()] = &state
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

func (m *streamsMap) UniStreamStates(config *ConnectionStateStoreConf) map[protocol.StreamID]*handover.UniStreamState {
	ss := make(map[protocol.StreamID]*handover.UniStreamState)
	for _, stream := range m.outgoingUniStreams.streams {
		s := &handover.UniStreamState{}
		stream.storeSendState(s, m.perspective, config)
		ss[stream.StreamID()] = s
	}
	for _, entry := range m.incomingUniStreams.streams {
		stream := entry.stream
		s := &handover.UniStreamState{}
		stream.storeReceiveState(s, m.perspective, config)
		ss[stream.StreamID()] = s
	}
	return ss
}

func (m *streamsMap) StoreState(state *handover.State, config *ConnectionStateStoreConf) {
	state.UniStreams = m.UniStreamStates(config)
	state.BidiStreams = m.BidiStreamStates(config)
	sfp := state.FromPerspective(m.perspective)
	sfp.SetNextIncomingUniStream(m.incomingUniStreams.nextStreamToAccept.StreamID(protocol.StreamTypeUni, m.perspective.Opposite()))
	sfp.SetNextIncomingBidiStream(m.incomingBidiStreams.nextStreamToAccept.StreamID(protocol.StreamTypeBidi, m.perspective.Opposite()))
	sfp.SetNextOutgoingUniStream(m.outgoingUniStreams.nextStream.StreamID(protocol.StreamTypeUni, m.perspective))
	sfp.SetNextOutgoingBidiStream(m.outgoingBidiStreams.nextStream.StreamID(protocol.StreamTypeBidi, m.perspective))
	sfp.SetMaxIncomingUniStream(int64(m.incomingUniStreams.maxStream.StreamID(protocol.StreamTypeUni, m.perspective.Opposite())))
	sfp.SetMaxIncomingBidiStream(int64(m.incomingBidiStreams.maxStream.StreamID(protocol.StreamTypeUni, m.perspective.Opposite())))
	sfp.SetMaxOutgoingUniStream(int64(m.outgoingUniStreams.maxStream.StreamID(protocol.StreamTypeBidi, m.perspective)))
	sfp.SetMaxOutgoingBidiStream(int64(m.outgoingBidiStreams.maxStream.StreamID(protocol.StreamTypeBidi, m.perspective)))
}

func (m *streamsMap) Restore(state *handover.State) (*RestoredStreams, error) {
	restoredStreams := &RestoredStreams{
		BidiStreams:    make(map[StreamID]Stream, 0),
		ReceiveStreams: make(map[StreamID]ReceiveStream, 0),
		SendStreams:    make(map[StreamID]SendStream, 0),
	}
	for streamID, streamState := range state.BidiStreams {
		stream, err := m.RestoreBidiStream(streamID, streamState)
		if err != nil {
			return nil, err
		}
		restoredStreams.BidiStreams[stream.StreamID()] = stream
	}
	for streamID, streamState := range state.UniStreams {
		if streamID.InitiatedBy() == m.perspective {
			stream, err := m.RestoreSendStream(streamID, streamState)
			if err != nil {
				return nil, err
			}
			restoredStreams.SendStreams[stream.StreamID()] = stream
		} else {
			stream, err := m.RestoreReceiveStream(streamID, streamState)
			if err != nil {
				return nil, err
			}
			restoredStreams.ReceiveStreams[stream.StreamID()] = stream
		}
	}
	sfp := state.FromPerspective(m.perspective)
	m.outgoingUniStreams.nextStream = sfp.NextOutgoingUniStream().StreamNum()
	m.outgoingBidiStreams.nextStream = sfp.NextOutgoingBidiStream().StreamNum()
	m.outgoingUniStreams.maxStream = protocol.StreamID(sfp.MaxOutgoingUniStream()).StreamNum()
	m.outgoingBidiStreams.maxStream = protocol.StreamID(sfp.MaxOutgoingBidiStream()).StreamNum()
	m.incomingUniStreams.nextStreamToAccept = sfp.NextIncomingUniStream().StreamNum()
	m.incomingUniStreams.nextStreamToOpen = m.incomingUniStreams.nextStreamToAccept
	m.incomingUniStreams.maxStream = protocol.StreamID(sfp.MaxIncomingUniStream()).StreamNum()
	m.incomingBidiStreams.nextStreamToAccept = sfp.NextIncomingBidiStream().StreamNum()
	m.incomingBidiStreams.nextStreamToOpen = m.incomingBidiStreams.nextStreamToAccept
	m.incomingBidiStreams.maxStream = protocol.StreamID(sfp.MaxIncomingBidiStream()).StreamNum()

	return restoredStreams, nil
}
