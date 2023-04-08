package quic

import (
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type queuedDatagram struct {
	awaitDequeue bool
	frame        *wire.DatagramFrame
}

type datagramQueue struct {
	sendQueue chan *queuedDatagram
	nextFrame *queuedDatagram
	rcvQueue  chan []byte

	closeErr error
	closed   chan struct{}

	hasData func()

	dequeued chan struct{}

	logger utils.Logger
}

func newDatagramQueue(hasData func(), logger utils.Logger) *datagramQueue {
	return &datagramQueue{
		hasData:   hasData,
		sendQueue: make(chan *queuedDatagram, protocol.DatagramSendQueueLen),
		rcvQueue:  make(chan []byte, protocol.DatagramRcvQueueLen),
		dequeued:  make(chan struct{}),
		closed:    make(chan struct{}),
		logger:    logger,
	}
}

// AddWithoutWaitForDequeue queues a new DATAGRAM frame for sending.
// Does not block until dequeue.
// Might block because send queue is full.
func (h *datagramQueue) AddWithoutWaitForDequeue(f *wire.DatagramFrame) error {
	select {
	case h.sendQueue <- &queuedDatagram{awaitDequeue: false, frame: f}:
		h.hasData()
		return nil
	case <-h.closed:
		return h.closeErr
	}
}

// AddAndWait queues a new DATAGRAM frame for sending.
// It blocks until the frame has been dequeued.
func (h *datagramQueue) AddAndWait(f *wire.DatagramFrame) error {
	select {
	case h.sendQueue <- &queuedDatagram{awaitDequeue: true, frame: f}:
		h.hasData()
	case <-h.closed:
		return h.closeErr
	}

	select {
	case <-h.dequeued:
		return nil
	case <-h.closed:
		return h.closeErr
	}
}

// Peek gets the next DATAGRAM frame for sending.
// If actually sent out, Pop needs to be called before the next call to Peek.
func (h *datagramQueue) Peek() *wire.DatagramFrame {
	if h.nextFrame != nil {
		return h.nextFrame.frame
	}
	select {
	case h.nextFrame = <-h.sendQueue:
		if h.nextFrame.awaitDequeue {
			h.dequeued <- struct{}{}
		}
	default:
		return nil
	}
	return h.nextFrame.frame
}

func (h *datagramQueue) Pop() {
	if h.nextFrame == nil {
		panic("datagramQueue BUG: Pop called for nil frame")
	}
	h.nextFrame = nil
}

// HandleDatagramFrame handles a received DATAGRAM frame.
func (h *datagramQueue) HandleDatagramFrame(f *wire.DatagramFrame) {
	data := make([]byte, len(f.Data))
	copy(data, f.Data)
	select {
	case h.rcvQueue <- data:
	default:
		h.logger.Debugf("Discarding DATAGRAM frame (%d bytes payload)", len(f.Data))
	}
}

// Receive gets a received DATAGRAM frame.
func (h *datagramQueue) Receive() ([]byte, error) {
	select {
	case data := <-h.rcvQueue:
		return data, nil
	case <-h.closed:
		return nil, h.closeErr
	}
}

func (h *datagramQueue) CloseWithError(e error) {
	h.closeErr = e
	close(h.closed)
}
