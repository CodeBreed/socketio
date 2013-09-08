// Package socketio implements a client for SocketIO protocol
// as specified in https://github.com/LearnBoost/socket.io-spec
package socketio

import (
	"time"
)

type Socket struct {
	Url       string
	Session   *Session
	Transport Transport
}

func Dial(url string, channel string, query string) (*Socket, error) {
	session, err := NewSession(url)
	if err != nil {
		return nil, err
	}

	transport, err := NewWSTransport(session, url)
	if err != nil {
		return nil, err
	}

	// Connect
	endpoint := NewEndpoint(channel, query)
	connectMsg := NewConnect(endpoint)
	transport.send(connectMsg)

	// Heartbeat goroutine
	go func() {
		heartbeatMsg := NewHeartbeat()
		for {
			time.Sleep(time.Duration(session.HeartbeatTimeout-1) * time.Second)
			_ = transport.send(heartbeatMsg)
		}
	}()

	return &Socket{url, session, transport}, nil
}

func (socket *Socket) Receive() (*Message, error) {
	return socket.Transport.receive()
}

func (socket *Socket) Send(msg *Message) error {
	return socket.Transport.send(msg)
}
