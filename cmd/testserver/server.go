package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"nhooyr.io/websocket"
)

// echoServer is the WebSocket echo server implementation.
// It ensures the client speaks the echo subprotocol and
// only allows one message every 100ms with a 10 message burst.
type echoServer struct {

	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		fmt.Printf("loopy\n")
		err = echo(r.Context(), c, l)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			s.logf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	fmt.Printf("before wait\n")
	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("after wait\n")
	_, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	// get byte array, unmarshall
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	s := string(buf.Bytes())
	/*
		var state common.GameState
		err = json.Unmarshal(buf.Bytes(), &state)
		if err != nil {
			fmt.Printf("unable to unmarshal gamestate %s\n", err.Error())
		}

		fmt.Printf("returned state %v\n", state) */
	fmt.Printf("data is %s\n", s)
	return err
}
