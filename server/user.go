package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	Conn   net.Conn
	C      chan string
	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		Conn:   conn,
		C:      make(chan string),
		Server: server,
	}
	go user.ListenMessage()
	return user
}

func (u *User) OnLine() {
	u.Server.mapLock.Lock()
	u.Server.Users[u.Name] = u
	u.Server.mapLock.Unlock()

	u.Server.BroadCast(u, "已上线")
}

func (u *User) OffLine() {
	u.Server.mapLock.Lock()
	delete(u.Server.Users, u.Name)
	u.Server.mapLock.Unlock()

	u.Server.BroadCast(u, "已下线")
}
func (u *User) SendMsg(msg string) {
	u.Conn.Write([]byte(msg))
}

func (u *User) DoMsg(msg string) {
	if msg == "who" {
		u.Server.mapLock.Lock()
		for _, user := range u.Server.Users {
			sendMsg := "[" + user.Addr + "]" + user.Name + ": is online...\n"
			u.SendMsg(sendMsg)
		}
		u.Server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := u.Server.Users[newName]; ok {
			u.SendMsg("用户名已存在\n")
		} else {
			u.Server.mapLock.Lock()
			delete(u.Server.Users, u.Name)
			u.Server.Users[newName] = u
			u.Server.mapLock.Unlock()

			u.Name = newName
			u.SendMsg("你已更改用户名为" + newName + "...\n")
		}
	} else if len(msg) > 3 && msg[:3] == "to|" {
		remoteUser := strings.Split(msg, "|")[1]
		if remoteUser == "" {
			u.SendMsg("消息格式不正确\n")
		}
		if _, ok := u.Server.Users[remoteUser]; !ok {
			u.SendMsg("用户不存在\n")
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			u.SendMsg("消息内容不能为空\n")
		}
		remoteUs := u.Server.Users[remoteUser]
		remoteUs.SendMsg(u.Name + "对你说:" + content + "\n")
	} else {
		u.Server.BroadCast(u, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg))
	}
}
