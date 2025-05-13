package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip      string
	Port    int
	Users   map[string]*User
	mapLock sync.RWMutex

	// Message channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:      ip,
		Port:    port,
		Users:   make(map[string]*User),
		Message: make(chan string),
	}
}

func (s *Server) Handler(conn net.Conn) {
	defer conn.Close()
	//fmt.Println("连接成功")
	user := NewUser(conn, s)
	user.OnLine()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMsg(msg)
			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:
		case <-time.After(120 * time.Second):
			user.SendMsg("你被踢了\n")
			user.OffLine()
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg + "\n"
	s.Message <- sendMsg
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.Users {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Start() {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("listen err:", err)
		return
	}
	defer listen.Close()
	go s.ListenMessage()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err:", err)
			continue
		}
		go s.Handler(conn)
	}
}
