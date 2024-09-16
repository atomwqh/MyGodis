package tcp

import (
	"bufio"
	"math/rand"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestListenAndServe(t *testing.T) {
	var err error
	closeChan := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
		return
	}
	addr := listener.Addr().String()
	go ListenAndServe(listener, NewEchoHandler(), closeChan)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		val := strconv.Itoa(rand.Int())
		_, err = conn.Write([]byte(val + "\n"))
		if err != nil {
			t.Error(err)
			return
		}
		bufReader := bufio.NewReader(conn)
		line, _, err := bufReader.ReadLine()
		if err != nil {
			t.Error(err)
			return
		}
		if string(line) != val {
			t.Error(string(line))
			return
		}
	}
	_ = conn.Close()
	for i := 0; i < 5; i++ {
		_, _ = net.Dial("tcp", addr)
	}
	closeChan <- struct{}{}
	time.Sleep(time.Second)
}

func TestClientCounter(t *testing.T) {
	var err error
	closeChan := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
		return
	}
	addr := listener.Addr().String()
	go ListenAndServe(listener, NewEchoHandler(), closeChan)

	sleepUntil := time.Now().Add(1 * time.Second)
	subtime := func() time.Duration {
		return sleepUntil.Sub(time.Now())
	}
	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				t.Error(err.Error())
			}
			defer conn.Close()
			time.Sleep(subtime())
		}()
		time.Sleep(5 * time.Microsecond)
	}
	time.Sleep(time.Second * 3)
	if ClientCounter != 0 {
		t.Errorf("Client Counter error :%d", ClientCounter)
	}
}
