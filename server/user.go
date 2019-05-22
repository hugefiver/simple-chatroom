package main

import (
	"fmt"
	"net"
	"time"
)

type Msg struct {
	from string
	text string
}

type User struct {
	Name string
	conn net.Conn
	sendMsgsChan chan *Msg
	commonMsgChan chan *Msg
	disconnect chan struct{}
}

func (u *User) SendMsg(msg *Msg) {
	if msg.from != u.Name {
		u.sendMsgsChan <- msg
	}
}

func (u *User) Handle() {
	defer server.Leaves(u.Name)

	go func() {
		for {
			select {
			case <-u.disconnect:
				return
			case msg := <-u.sendMsgsChan:
				u.conn.Write([]byte(fmt.Sprintf("[%s]<%s> %s", time.Now().Format("15:04:05"), msg.from, msg.text)))
			}
		}
	}()

	for {
		b := make([]byte, 1024)
		l, err := u.conn.Read(b)
		if err != nil {
			u.disconnect <- struct {}{}
			return
		}

		u.commonMsgChan <- &Msg{
			from: u.Name,
			text: string(b[:l]),
		}
	}
}

func NewUser(name string, conn net.Conn, commonChan chan *Msg) *User {
	return &User{
		Name:name,
		conn:conn,
		commonMsgChan:commonChan,
		sendMsgsChan:make(chan *Msg, 1024),
		disconnect:make(chan struct{}, 1),
	}
}