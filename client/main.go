package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

var (
	addr = os.Args[1]
)

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Can't connect to", addr)
	}
	defer conn.Close()

	var name string
	fmt.Print("请输入用户名(20字节以内): ")

	if _, err := fmt.Scanf("%s", &name); err != nil {
		log.Fatalln("读取用户名错误", name)
	}

	conn.Write([]byte(name))
	fmt.Println("已登入.输入想说的话后回车即可.")

	go func() {
		buff := make([]byte, 1024)
		for{
			n, err := conn.Read(buff)
			if err != nil {
				log.Fatalln("[连接出错]", err)
			} else {
				fmt.Println(string(buff[:n]))
			}
		}
	}()

	for {
		buff := make([]byte, 1024)
		_, _ = fmt.Scanln(&buff)
		conn.Write(buff)
	}
}
