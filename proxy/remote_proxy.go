package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// RemoteProxy struct
type RemoteProxy struct {
	from string
	to   string
	done chan struct{}
}

// NewRemoteProxy ctor
func NewRemoteProxy(from, to string) *RemoteProxy {
	return &RemoteProxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
}

// Start the proxy
func (p *RemoteProxy) Start() error {
	fmt.Printf("Starting proxy...\n")
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	p.run(listener)
	return nil
}

// Stop the proxy
func (p *RemoteProxy) Stop() {
	fmt.Printf("stopping proxy...\n")
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func (p *RemoteProxy) run(listener net.Listener) {
	for {
		select {
		case <-p.done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.handle(connection)
			} else {
				fmt.Printf("error accepting conn: %s\n", err.Error())
			}
		}
	}
}

func (p *RemoteProxy) handle(connection net.Conn) {
	fmt.Printf(
		"proxying from %s to %s\n",
		connection.RemoteAddr().String(),
		p.to,
	)
	defer fmt.Printf("done handling %p\n", connection)
	defer connection.Close()
	auth := make([]byte, 1024)
	// receive authentication bytes
	_, err := connection.Read(auth)
	if err != nil {
		fmt.Printf("failed to read authentication\n")
	} else {
		remote, err := net.Dial("tcp", p.to)
		if err != nil {
			fmt.Printf("error dialing remote host: %s\n", err.Error())
			return
		}
		defer remote.Close()
		fmt.Printf("remote part authenticated as: %s\n", string(auth))
		// send authentication confirmation after connecting to upstream
		connection.Write([]byte("auth ok, continue"))
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go p.copy(remote, connection, wg)
		go p.copy(connection, remote, wg)
		wg.Wait()
	}
}

func (p *RemoteProxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case <-p.done:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			fmt.Printf("error copying: %s\n", err.Error())
			p.Stop()
			return
		}
	}
}
