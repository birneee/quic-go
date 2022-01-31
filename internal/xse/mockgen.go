package xse

//go:generate sh -c "mockgen -package xse -destination mock_receive_stream_test.go github.com/lucas-clemente/quic-go/internal/xse ReceiveStream && goimports -w interface.go"
//go:generate sh -c "mockgen -package xse -destination mock_send_stream_test.go github.com/lucas-clemente/quic-go/internal/xse SendStream && goimports -w interface.go"
