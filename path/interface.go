package path

import "net"

type PathManager interface {
	IgnoreReceiveFrom(addr net.Addr)
	IgnoreMigrateTo(addr net.Addr)
	IgnoreSendTo(addr net.Addr)
	IsIgnoreReceiveFrom(addr net.Addr) bool
	IsIgnoreMigrateTo(addr net.Addr) bool
	IsIgnoreSendTo(addr net.Addr) bool
}
