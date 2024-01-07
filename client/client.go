package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	client.Name = conn.LocalAddr().String()
	return client
}

var server_ip string
var server_port int

func init() {
	flag.StringVar(&server_ip, "ip", "127.0.0.1", "默认ip地址是127.0.0.1")
	flag.IntVar(&server_port, "port", 8888, "默认端口是8888")
	flag.Parse()
}
func (this *Client) menu() bool {
	var flag int
	fmt.Println("当前用户：" + this.Name)
	fmt.Println("1.公聊模式\n2.私聊模式\n3.更新用户名\n0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范围内的数字")
		return false
	}

}
func (this *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名：")
	fmt.Scanln(&this.Name)
	for this.Name == "" {
		fmt.Println(">>>>>用户名不能为空！！")
		fmt.Println(">>>>>请输入用户名：")
		fmt.Scanln(&this.Name)
	}
	sendMsg := "rename|" + this.Name
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}
func (this *Client) Run() {
	for this.flag != 0 {
		for !this.menu() {

		}
		switch this.flag {
		case 1:
			// fmt.Println("公聊模式选择...")
			this.PublicChat()
			break
		case 2:
			// fmt.Println("私聊模式选择...")
			this.PrivateChat()
			break
		case 3:
			// fmt.Println("更新用户名选择...")
			this.UpdateName()
			break
		case 0:
			// fmt.Println("退出选择...")
			os.Exit(0)
			break
		}
	}
}
func (this *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>>请输入聊天内容,exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			_, err := this.conn.Write([]byte(chatMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)
	}
}
func (this *Client) SelectUsers() {
	sendMsg := "who"
	_, err := this.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}
func (this *Client) PrivateChat() {
	var remoteName, chatMsg string

	this.SelectUsers()
	fmt.Println(">>>>>请输入聊天对象【用户名】,exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容,exit退出")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg
				_, err := this.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>请输入聊天内容,exit推出")
			fmt.Scanln(&chatMsg)
		}
		this.SelectUsers()
		fmt.Println(">>>>>请输入聊天对象【用户名】,exit退出")
		fmt.Scanln(&remoteName)
	}

}
func (this *Client) DealResponse() {
	io.Copy(os.Stdout, this.conn)
}
func main() {
	client := NewClient(server_ip, server_port)
	if client == nil {
		fmt.Println(">>>>>连接服务器失败...")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>连接服务器成功...")
	client.Run()
	select {}
}
