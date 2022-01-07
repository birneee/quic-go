package mocks

//go:generate sh -c "mockgen -package mocks -destination receive_stream.go github.com/lucas-clemente/quic-go/internal/xse ReceiveStream && goimports -w receive_stream.go"
//go:generate sh -c "mockgen -package mocks -destination send_stream.go github.com/lucas-clemente/quic-go/internal/xse SendStream && goimports -w send_stream.go"
