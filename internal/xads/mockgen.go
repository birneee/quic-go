package xads

//go:generate sh -c "mockgen -package xads -destination mock_receive_stream_test.go github.com/quic-go/quic-go/internal/xads ReceiveStream && goimports -w interface.go"
//go:generate sh -c "mockgen -package xads -destination mock_send_stream_test.go github.com/quic-go/quic-go/internal/xads SendStream && goimports -w interface.go"
