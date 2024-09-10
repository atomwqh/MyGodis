package tcp

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

func ListenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(fmt.Sprintf("listen tcp err:%s", err))
	}
	defer listener.Close()
	log.Println(fmt.Sprintf("listen tcp success, bind : %s", address))

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(fmt.Sprintf("accept err:%v", err))
		}
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				log.Println("connection closed")
			} else {
				log.Println(err)
			}
			return
		}
		b := []byte(msg)
		conn.Write(b)
	}
}

func main() {
	ListenAndServe(":8000")
}
