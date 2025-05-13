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
	Conn       net.Conn
	Flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Flag:       999,
	}
	conn, err := net.Dial("tcp", net.JoinHostPort(serverIp, fmt.Sprintf("%d", serverPort)))
	if err != nil {
		fmt.Println("net dial err:", err)
		return nil
	}
	client.Conn = conn
	return client
}

func (c *Client) Menu() bool {
	fmt.Println("1. 公聊模式")
	fmt.Println("2. 私聊模式")
	fmt.Println("3. 修改用户名")
	fmt.Println("0. 退出")
	//fmt.Print("请输入选项：")

	var flag int
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		c.Flag = flag
		return true
	} else {
		fmt.Println("请输入合法数字范围(0-3)>>>>>>>")
		return false
	}
}

func (c *Client) UpdateName() bool {
	fmt.Print(">>>>>>请输入用户名：")
	var name string
	fmt.Scanln(&name)
	SendMsg := "rename|" + name + "\n"
	if _, err := c.Conn.Write([]byte(SendMsg)); err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (c *Client) run() {
	for c.Flag != 0 {
		for !c.Menu() {
		}
		switch c.Flag {
		case 1:
			fmt.Println("公聊模式选择...")
		case 2:
			fmt.Println("私聊模式选择...")
		case 3:
			fmt.Println("修改用户名选择...")
			c.UpdateName()
		}
	}
}

var serverIp string
var serverPort int

func init() {
	//fmt.Println(">>>>>>>>>>>>>>>>客户端启动成功>>>>>>>>>>>>>>>")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器ip地址设置(默认127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "服务器ip地址设置(默认8888)")
}

func (c *Client) DelResponse() {
	io.Copy(os.Stdout, c.Conn)
}
func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>>>>>>>>服务器连接失败>>>>>>>>>>>>>>>")
		return
	}
	fmt.Println(">>>>>>>>>>>>>>>>服务器连接成功>>>>>>>>>>>>>>>")
	go client.DelResponse()
	client.run()
}
