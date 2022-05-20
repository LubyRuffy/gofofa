package gofofa

import (
	"compress/gzip"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fetchHander = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/gzip.json":
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			gz.Write([]byte(`{"text":"hello world"}`))
			return
		case "/api/v1/contentLengthError.json":
			w.Header().Set("Content-Length", "aaa")
			return
		}
	}
)

type tcpTestServer struct {
	listener net.Listener
}

func (ts *tcpTestServer) URL() string {
	return "http://" + ts.listener.Addr().String()
}

func (ts *tcpTestServer) Close() {
	ts.listener.Close()
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	return l
}

func newTcpTestServer(handler func(conn net.Conn, data []byte) error) *tcpTestServer {
	ts := &tcpTestServer{
		listener: newLocalListener(),
	}

	go func() {
		data := make([]byte, 1024)
		var err error
		var n int
		var conn net.Conn

		for {
			conn, err = ts.listener.Accept()
			if err != nil {
				break
			}

			n, err = conn.Read(data)
			if err != nil {
				break
			}

			err = handler(conn, data[:n])
			if err != nil {
				break
			}
			conn.Close()
		}
	}()

	return ts
}

func TestClient_Fetch(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:55")
	assert.Error(t, err)

	ts := httptest.NewServer(http.HandlerFunc(fetchHander))
	defer ts.Close()

	cli := &Client{
		Server:     ts.URL,
		APIVersion: "v1",
		httpClient: &http.Client{},
	}

	// 解析异常
	var a map[string]interface{}
	err = cli.Fetch("", nil, &a)
	assert.Error(t, err)

	// gzip
	err = cli.Fetch("gzip.json", nil, &a)
	assert.Nil(t, err)
	assert.Equal(t, "hello world", a["text"].(string))

	// content Length Error
	err = cli.Fetch("contentLengthError.json", nil, &a)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")

	// read all error 需要构造一个恶意服务器
	s := newTcpTestServer(func(conn net.Conn, data []byte) error {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 10\r\nConnection: close\r\n\r\na"))
		return err
	})
	defer s.Close()
	cli = &Client{
		Server:     s.URL(),
		APIVersion: "v1",
		httpClient: &http.Client{},
	}
	err = cli.Fetch("/", nil, &a)
	assert.Contains(t, err.Error(), "unexpected EOF")

	// 构造content length 的atoi 异常， 需要构造一个恶意服务器
	s1 := newTcpTestServer(func(conn net.Conn, data []byte) error {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: -\r\nConnection: close\r\n\r\na"))
		return err
	})
	defer s1.Close()
	cli = &Client{
		Server:     s1.URL(),
		APIVersion: "v1",
		httpClient: &http.Client{},
	}
	err = cli.Fetch("/", nil, &a)
	assert.Contains(t, err.Error(), "bad Content-Length")

	// 构造错误的gzip， 需要构造一个恶意服务器
	s2 := newTcpTestServer(func(conn net.Conn, data []byte) error {
		_, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 1\r\nContent-Encoding: gzip\r\nConnection: close\r\n\r\na"))
		return err
	})
	defer s2.Close()
	cli = &Client{
		Server:     s2.URL(),
		APIVersion: "v1",
		httpClient: &http.Client{},
	}
	err = cli.Fetch("/", nil, &a)
	assert.Contains(t, err.Error(), "unexpected EOF")
}
