package cli

/*
import (
	"net"

	logger "github.com/sirupsen/logrus"
)

// NewSocket() (Socket, err) is defined in the various platform-specific socket_*.go files.
type Socket interface {
	BindToSocket() (net.Listener, error)
	DialSocket() (net.Conn, error)
}

type SocketInfo struct {
	log       logger.Logger
	bindFile  string
	dialFiles []string
	testOwner bool //nolint
}

func (s SocketInfo) GetBindFile() string {
	return s.bindFile
}

func (s SocketInfo) GetDialFiles() []string {
	return s.dialFiles
}

type SocketWrapper struct {
	Conn net.Conn
	// Transporter rpc.Transporter
	Err error
}
*/