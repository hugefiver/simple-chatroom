package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	commonMsgsChan chan *Msg
	NoticeChan chan string
	users map[string] *User
	bind string
}

func (s *Server) Leaves(user string) {
	s.NoticeChan <- fmt.Sprintf("[%s Leaves]\n", user)
	log.Printf("[%s Leaves]\n", user)
	delete(s.users, user)
}

func (s *Server) Handle() {
	server, err := net.Listen("tcp", s.bind)
	if err != nil {
		log.Fatal("Failed to listen at " + s.bind, err)
	}
	defer server.Close()

	go func() {
		for {
			var msg *Msg
			select {
			case msg = <-s.commonMsgsChan:
			case notice := <-s.NoticeChan:
				msg = &Msg{
					from:"Admin",
					text:notice,
				}
			}

			for _, v := range s.users  {
				v.SendMsg(msg)
			}
		}
	}()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("[Connect Error]", err)
		}
		
		go func(conn net.Conn) {
			defer conn.Close()
			name := make([]byte, 22)
			l, err := conn.Read(name)
			if err != nil {
				log.Println("[Regular Error]", err)
			}
			
			n := strings.TrimSpace(string(name[:l]))
			if _, ok := s.users[n]; ok != true && n != "Admin" {
				s.NoticeChan <- fmt.Sprintf("[%s] comes in", n)
				log.Println("[User Regular]", n)

				u := NewUser(n, conn, s.commonMsgsChan)
				s.users[n] = u
				u.Handle()
			} else {
				conn.Write([]byte("User Name Already Used.\n"))
				log.Println("User Name Already Used:", n)
			}
		}(conn)
	}
}

func NewServer(bind string) Server {
	return Server {
		commonMsgsChan:make(chan *Msg, 1024),
		NoticeChan:make(chan string, 1024),
		users: make(map[string] *User),
		bind:bind,
	}
}
