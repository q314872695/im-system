package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port}
	server.OnlineMap = make(map[string]*User)
	server.Message = make(chan string, 10)
	return server
}
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		// fmt.Println("server监听消息:", msg)
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
func (this *Server) BroadCast(user *User, msg string) {

	sendMsg := fmt.Sprintf("[%v]%v:%v\n", user.Addr, user.Name, msg)
	fmt.Println("server广播消息:", sendMsg)
	this.Message <- sendMsg
}
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功...")

	user := NewUser(conn, this)
	// fmt.Println("123-----------Online")
	user.Online()
	// fmt.Println("456-----------Online")
	isLive := make(chan bool)
	go func() {
		// fmt.Println("456-----------")
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("ConnRead err:", err)
				return
			}

			msg := string(buf[:n])
			// fmt.Printf("%v:%v\n", msg, msg == "who")
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	// fmt.Println("789-----------")
	for {
		select {
		case <-isLive:
		case <-time.After(300 * time.Second):
			user.SendMsg("你被踢了！！！\n")
			close(user.C)
			conn.Close()
			return
		}
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	defer listener.Close()
	go this.ListenMessager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go this.Handler(conn)
	}
}
