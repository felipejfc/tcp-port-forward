package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// LocalProxy struct
type LocalProxy struct {
	from string
	to   string
	done chan struct{}
}

// NewLocalProxy ctor
func NewLocalProxy(from, to string) *LocalProxy {
	return &LocalProxy{
		from: from,
		to:   to,
		done: make(chan struct{}),
	}
}

// Start the proxy
func (p *LocalProxy) Start() error {
	fmt.Printf("Starting proxy...\n")
	listener, err := net.Listen("tcp", p.from)
	if err != nil {
		return err
	}
	p.run(listener)
	return nil
}

// Stop the proxy
func (p *LocalProxy) Stop() {
	fmt.Printf("stopping proxy...\n")
	if p.done == nil {
		return
	}
	close(p.done)
	p.done = nil
}

func (p *LocalProxy) run(listener net.Listener) {
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

func (p *LocalProxy) handle(connection net.Conn) {
	fmt.Printf(
		"proxying from %s to %s\n",
		connection.RemoteAddr().String(),
		p.to,
	)
	defer fmt.Printf("done handling %p\n", connection)
	defer connection.Close()
	remote, err := net.Dial("tcp", p.to)
	if err != nil {
		fmt.Printf("error dialing remote host: %s\n", err.Error())
		return
	}
	defer remote.Close()
	// connect to remote proxy and send authentication packet
	_, err = remote.Write([]byte("felipe-cavalcanti-token"))
	if err != nil {
		fmt.Printf("failed to send authentication\n")
	}
	authRes := make([]byte, 1024)
	// wait for remote proxy authentication confirmation
	_, err = remote.Read(authRes)
	if err != nil {
		fmt.Printf("authentication error: %s\n", err.Error())
	} else {
		fmt.Printf("success: %s\n", string(authRes))
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go p.copy(remote, connection, wg)
	go p.copy(connection, remote, wg)
	wg.Wait()
}

func (p *LocalProxy) copy(from, to net.Conn, wg *sync.WaitGroup) {
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
