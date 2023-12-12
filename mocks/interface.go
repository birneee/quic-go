package mocks

import (
	mocklogging "github.com/quic-go/quic-go/internal/mocks/logging"
	"github.com/quic-go/quic-go/logging"
	"go.uber.org/mock/gomock"
)

type MockConnectionTracer = mocklogging.MockConnectionTracer

func NewMockConnectionTracer(ctrl *gomock.Controller) (*logging.ConnectionTracer, *MockConnectionTracer) {
	return mocklogging.NewMockConnectionTracer(ctrl)
}
