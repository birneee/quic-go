package quic

import (
	"context"
	"errors"
	"fmt"
	"github.com/quic-go/quic-go/handover"
	"github.com/quic-go/quic-go/internal/ackhandler"
	"github.com/quic-go/quic-go/qstate"
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
	ctx         context.Context // not used for cancellations, but carries the values associated with the connection
	perspective protocol.Perspective

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
}

var _ streamManager = &streamsMap{}

func newStreamsMap(
	ctx context.Context,
	sender streamSender,
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController,
	maxIncomingBidiStreams uint64,
	maxIncomingUniStreams uint64,
	perspective protocol.Perspective,
) streamManager {
	m := &streamsMap{
		ctx:                    ctx,
		perspective:            perspective,
		newFlowController:      newFlowController,
		maxIncomingBidiStreams: maxIncomingBidiStreams,
		maxIncomingUniStreams:  maxIncomingUniStreams,
		sender:                 sender,
	}
	m.initMaps()
	return m
}

func (m *streamsMap) initMaps() {
	m.outgoingBidiStreams = newOutgoingStreamsMap(
		protocol.StreamTypeBidi,
		func(num protocol.StreamNum) streamI {
			id := num.StreamID(protocol.StreamTypeBidi, m.perspective)
			return newStream(m.ctx, id, m.sender, m.newFlowController(id))
		},
		m.sender.queueControlFrame,
	)
	m.incomingBidiStreams = newIncomingStreamsMap(
		protocol.StreamTypeBidi,
		func(num protocol.StreamNum) streamI {
			id := num.StreamID(protocol.StreamTypeBidi, m.perspective.Opposite())
			return newStream(m.ctx, id, m.sender, m.newFlowController(id))
		},
		m.maxIncomingBidiStreams,
		m.sender.queueControlFrame,
	)
	m.outgoingUniStreams = newOutgoingStreamsMap(
		protocol.StreamTypeUni,
		func(num protocol.StreamNum) sendStreamI {
			id := num.StreamID(protocol.StreamTypeUni, m.perspective)
			return newSendStream(m.ctx, id, m.sender, m.newFlowController(id))
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

func (m *streamsMap) RestoreBidiStream(streamID StreamID, state *qstate.Stream) (Stream, error) {
	var stream Stream
	var err error
	if streamID.InitiatedBy() == m.perspective {
		stream, err = RestoreOutgoingBidiStream(m.outgoingBidiStreams, streamID.StreamNum(), state)
	} else {
		stream, err = RestoreIncomingBidiStream(m.incomingBidiStreams, streamID.StreamNum(), state)
	}
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (m *streamsMap) RestoreSendStream(streamID StreamID, state *qstate.Stream) (SendStream, error) {
	return RestoreOutgoingUniStream(m.outgoingUniStreams, streamID.StreamNum(), state)
}

func (m *streamsMap) RestoreReceiveStream(streamID StreamID, state *qstate.Stream) (ReceiveStream, error) {
	return RestoreIncomingUniStream(m.incomingUniStreams, streamID.StreamNum(), state)
}

func streamIToBidiStreamState(s streamI, sph ackhandler.SentPacketHandler, config *handover.ConnectionStateStoreConf) qstate.Stream {
	ss := qstate.Stream{
		StreamID: int64(s.StreamID()),
	}
	s.storeReceiveState(&ss, config)
	s.storeSendState(&ss, sph, config)
	return ss
}

func (m *streamsMap) AppendBidiStreamStates(streamStates []qstate.Stream, sph ackhandler.SentPacketHandler, config *handover.ConnectionStateStoreConf) []qstate.Stream {
	for _, stream := range m.outgoingBidiStreams.streams {
		streamStates = append(streamStates, streamIToBidiStreamState(stream, sph, config))
	}
	for _, entry := range m.incomingBidiStreams.streams {
		streamStates = append(streamStates, streamIToBidiStreamState(entry.stream, sph, config))
	}
	return streamStates
}

func (m *streamsMap) OpenedBidiStream(id StreamID) (Stream, error) {
	if id.InitiatedBy() == m.perspective {
		return m.outgoingBidiStreams.GetStream(id.StreamNum())
	} else {
		return m.incomingBidiStreams.GetStream(id.StreamNum())
	}
}

func (m *streamsMap) AppendUniStreamStates(streamStates []qstate.Stream, sph ackhandler.SentPacketHandler, config *handover.ConnectionStateStoreConf) []qstate.Stream {
	for _, stream := range m.outgoingUniStreams.streams {
		s := qstate.Stream{
			StreamID: int64(stream.StreamID()),
		}
		stream.storeSendState(&s, sph, config)
		streamStates = append(streamStates, s)
	}
	for _, entry := range m.incomingUniStreams.streams {
		s := qstate.Stream{
			StreamID: int64(entry.stream.StreamID()),
		}
		entry.stream.storeReceiveState(&s, config)
		streamStates = append(streamStates, s)
	}
	return streamStates
}

func (m *streamsMap) StoreState(state *qstate.Connection, sph ackhandler.SentPacketHandler, config *handover.ConnectionStateStoreConf) {
	state.Transport.Streams = m.AppendUniStreamStates(state.Transport.Streams, sph, config)
	state.Transport.Streams = m.AppendBidiStreamStates(state.Transport.Streams, sph, config)
	state.Transport.NextUnidirectionalStream = int64(m.outgoingUniStreams.nextStream.StreamID(protocol.StreamTypeUni, m.perspective))
	state.Transport.NextBidirectionalStream = int64(m.outgoingBidiStreams.nextStream.StreamID(protocol.StreamTypeBidi, m.perspective))
	state.Transport.MaxUnidirectionalStreams = int64(m.outgoingUniStreams.maxStream.StreamID(protocol.StreamTypeBidi, m.perspective))
	state.Transport.MaxBidirectionalStreams = int64(m.outgoingBidiStreams.maxStream.StreamID(protocol.StreamTypeBidi, m.perspective))
	state.Transport.RemoteNextUnidirectionalStream = int64(m.incomingUniStreams.nextStreamToAccept.StreamID(protocol.StreamTypeUni, m.perspective.Opposite()))
	state.Transport.RemoteNextBidirectionalStream = int64(m.incomingBidiStreams.nextStreamToAccept.StreamID(protocol.StreamTypeBidi, m.perspective.Opposite()))
	state.Transport.RemoteMaxUnidirectionalStreams = int64(m.incomingUniStreams.maxStream.StreamID(protocol.StreamTypeUni, m.perspective.Opposite()))
	state.Transport.RemoteMaxBidirectionalStreams = int64(m.incomingBidiStreams.maxStream.StreamID(protocol.StreamTypeUni, m.perspective.Opposite()))
}

func (m *streamsMap) restoreStreams(state *qstate.Connection) (*RestoredStreams, error) {
	restoredStreams := &RestoredStreams{
		BidiStreams:    make(map[StreamID]Stream, 0),
		ReceiveStreams: make(map[StreamID]ReceiveStream, 0),
		SendStreams:    make(map[StreamID]SendStream, 0),
	}

	for _, streamState := range state.Transport.Streams {
		streamID := protocol.StreamID(streamState.StreamID)
		switch streamID.Type() {
		case protocol.StreamTypeBidi:
			stream, err := m.RestoreBidiStream(streamID, &streamState)
			if err != nil {
				return nil, err
			}
			restoredStreams.BidiStreams[stream.StreamID()] = stream
		case protocol.StreamTypeUni:
			switch streamID.InitiatedBy() {
			case m.perspective:
				stream, err := m.RestoreSendStream(streamID, &streamState)
				if err != nil {
					return nil, err
				}
				restoredStreams.SendStreams[stream.StreamID()] = stream
			case m.perspective.Opposite():
				stream, err := m.RestoreReceiveStream(streamID, &streamState)
				if err != nil {
					return nil, err
				}
				restoredStreams.ReceiveStreams[stream.StreamID()] = stream
			default:
				panic("unexpected initiator")
			}
		default:
			panic("unexpected stream type")
		}
	}

	return restoredStreams, nil
}

func (m *streamsMap) SendStreams() []flowcontrol.SendStream {
	var sendStreams []flowcontrol.SendStream
	for _, s := range m.outgoingUniStreams.streams {
		sendStreams = append(sendStreams, s)
	}
	for _, s := range m.outgoingBidiStreams.streams {
		sendStreams = append(sendStreams, s)
	}
	for _, e := range m.incomingBidiStreams.streams {
		sendStreams = append(sendStreams, e.stream)
	}
	return sendStreams
}

// must be called after all field required by the newFlowController function are set
func restoreStreamMap(
	state *qstate.Connection,
	sender streamSender,
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController,
) (func() (*RestoredStreams, error), streamManager, error) {
	s := newStreamsMap(sender, newFlowController, uint64(state.Transport.MaxBidirectionalStreams), uint64(state.Transport.MaxUnidirectionalStreams), state.Transport.Perspective()).(*streamsMap)

	s.outgoingUniStreams.nextStream = protocol.StreamID(state.Transport.NextUnidirectionalStream).StreamNum()
	s.outgoingBidiStreams.nextStream = protocol.StreamID(state.Transport.NextBidirectionalStream).StreamNum()
	s.outgoingUniStreams.maxStream = protocol.StreamID(state.Transport.MaxUnidirectionalStreams).StreamNum()
	s.outgoingBidiStreams.maxStream = protocol.StreamID(state.Transport.MaxBidirectionalStreams).StreamNum()
	s.incomingUniStreams.nextStreamToAccept = protocol.StreamID(state.Transport.RemoteNextUnidirectionalStream).StreamNum()
	s.incomingUniStreams.nextStreamToOpen = s.incomingUniStreams.nextStreamToAccept
	s.incomingUniStreams.maxStream = protocol.StreamID(state.Transport.RemoteMaxUnidirectionalStreams).StreamNum()
	s.incomingBidiStreams.nextStreamToAccept = protocol.StreamID(state.Transport.RemoteNextBidirectionalStream).StreamNum()
	s.incomingBidiStreams.nextStreamToOpen = s.incomingBidiStreams.nextStreamToAccept
	s.incomingBidiStreams.maxStream = protocol.StreamID(state.Transport.RemoteMaxBidirectionalStreams).StreamNum()

	return func() (*RestoredStreams, error) {
		return s.restoreStreams(state)
	}, s, nil
}
