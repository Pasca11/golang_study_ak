package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// подключиться к серверу
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Connected to localhost:8000")
	// запустить горутину, которая будет читать все сообщения от сервера и выводить их в консоль
	go clientReader(conn)
	// читать сообщения от stdin и отправлять их на сервер
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		_, err := conn.Write([]byte(line + "\r\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// clientReader выводит на экран все сообщения от сервера
func clientReader(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println("From server", scanner.Text())
	}
}
