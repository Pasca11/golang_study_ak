package main

import (
	"bytes"
	"net"
	"testing"
	"time"
)

type mockConn struct {
	writeBytes []byte
	readBytes  []byte
}

func (m mockConn) Read(b []byte) (n int, err error) {
	m.readBytes = b
	return len(b), nil
}

func (m mockConn) Write(b []byte) (n int, err error) {
	m.writeBytes = b
	return len(b), nil
}

func (m mockConn) Close() error {
	return nil
}

func (m mockConn) LocalAddr() net.Addr {
	return nil
}

func (m mockConn) RemoteAddr() net.Addr {
	addr, _ := net.ResolveIPAddr("tcp", "127.0.0.1:8080")
	return addr
}

func (m mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestBroadcaster(t *testing.T) {
	ch1 := make(chan string, 2)
	ch2 := make(chan string, 2)
	client1 := client{ch: ch1}
	client2 := client{ch: ch2}

	go broadcaster()

	entering <- client1
	entering <- client2

	time.Sleep(time.Millisecond)
	messages <- "test message"

	msg := <-ch1
	if msg != "test message" {
		t.Errorf("expected 'test message', got %s", msg)
	}
	msg = <-ch2
	if msg != "test message" {
		t.Errorf("expected 'test message', got %s", msg)
	}
}

func TestHandleConnection(t *testing.T) {
	mock := &mockConn{
		readBytes: []byte("test message"),
	}

	go handleConn(mock)

	if !bytes.Equal(mock.readBytes, []byte("test message")) {
		t.Errorf("expected 'test message', got %s", mock.readBytes)
	}
}
