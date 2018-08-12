package common

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func ListenRpc(addr string, obj interface{}) error {
	defer CheckPanic()

	if obj == nil {
		return fmt.Errorf("rpc object is nil")
	}

	server := rpc.NewServer()
	server.Register(obj)

	if listener, err := net.Listen("tcp", addr); err != nil {
		return fmt.Errorf("rpc listen error : %v : %v", addr, err)
	} else {
		log.Println("rpc running @", addr)
		server.Accept(listener)
	}

	return fmt.Errorf("rpc server quit : %v", addr)
}

func ListenSocket(addr string, keepalive bool, reactiver func(net.Conn)) error {
	defer CheckPanic()

	if reactiver == nil {
		return fmt.Errorf("socket reactiver is nil")
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return fmt.Errorf("can't resolve addr : %v : %v", addr, err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("can't listen tcp : %v : %v", addr, err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("accept ", err)
			break
		}

		conn.SetNoDelay(true)
		conn.SetKeepAlive(keepalive)
		conn.SetLinger(-1)
		go reactiver(conn)
	}

	return fmt.Errorf("socket server quit : %v ", addr)
}

func ListenUdp(addr string, bufsize int, reactiver func(*net.UDPAddr, []byte) []byte) error {
	defer CheckPanic()

	if reactiver == nil {
		return fmt.Errorf("udp reactiver is nil")
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("can't resolve addr : %v : %v", addr, err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("can't listen udp : %v : %v", addr, err)
	}
	defer conn.Close()

	for {
		data := make([]byte, bufsize)
		nr, remote, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Println("read from udp ", remote, err)
			continue
		}

		if ret := reactiver(remote, data[:nr]); ret != nil && len(ret) > 0 {
			conn.WriteToUDP(ret, remote)
		}
	}

	return fmt.Errorf("udp server quit : %v ", addr)
}
