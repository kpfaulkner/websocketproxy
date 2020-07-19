package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/kpfaulkner/websocketproxy/common"
	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	//"io"
	"net/http"
	"nhooyr.io/websocket"
	"time"
)

type WebSocketPair struct {
	SourceWebSocket *websocket.Conn
	TargetWebSocket *websocket.Conn
}

// ProxyServer is the WebSocket proxy.... I think :)
type ProxyServer struct {

	// LUT between incoming and outgoing websockets
	websocketLUT map[string]WebSocketPair

	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

func NewProxyServer() ProxyServer {
	p := ProxyServer{}
  p.websocketLUT = make(map[string]WebSocketPair)

	return p
}

func (s ProxyServer) processEcho() {

}

func (s* ProxyServer) processWebSocketConnection( incomingCtx context.Context, incomingWS *websocket.Conn, remoteAddr string, host string ) {

	// check if connection already used.
	key := remoteAddr+":"+host
	var wsPair WebSocketPair
	var ok bool
	wsPair, ok = s.websocketLUT[key]
	if !ok {
	  // dont have entry
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)
		newConn, _, err := websocket.Dial(ctx, fmt.Sprintf("ws://%s",host), &websocket.DialOptions{})
		if err != nil {
			fmt.Printf("error dialing %s\n", err.Error())
			return
		}

		wsPair = WebSocketPair{SourceWebSocket: incomingWS, TargetWebSocket: newConn}
		s.websocketLUT[key] = wsPair
	}

	ctx := context.Background()
	quit := false
	// loop and purely read from source and write to target.
	// but in reality will need to modify due to WS's working in both directions.
	for !quit {

		typ, r, err := wsPair.SourceWebSocket.Reader(incomingCtx)
		if err != nil {
			log.Errorf("Unable to read from source websocket: %s\n", err.Error())
			return
		}

		// get byte array, unmarshall
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)

		s := string(buf.Bytes())
		fmt.Printf("data is %s\n", s)

		err = wsPair.TargetWebSocket.Write(ctx, typ, buf.Bytes())

		//_, err = io.Copy(w, r)
		if err != nil {
			log.Errorf("Unable to copy data between websockets: %s\n", err.Error())
			return
		}
	}
}

func (s ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		s.logf("%v", err)
		return
	}
	//defer c.Close(websocket.StatusInternalError, "the sky is falling")

	log.Debugf("remote addr is %s\n", r.RemoteAddr)
	log.Debugf("host addr is %s\n", r.Host)
	log.Debugf("websocket version is %s\n", r.Header.Get("Sec-Websocket-Version"))

	ctx := r.Context()
	// go and process this socket.
	s.processWebSocketConnection(ctx,c, r.RemoteAddr, r.Host)

}

// echo reads from the WebSocket connection and then writes
// the received message back to it.
// The entire function has 10s to complete.
func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	//fmt.Printf("in echo\n")
	//err := l.Wait(ctx)
	//if err != nil {
	//	return err
	//}

	_, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	// get byte array, unmarshall
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

	var state common.GameState
	err = json.Unmarshal(buf.Bytes(), &state)
	if err != nil {
		fmt.Printf("unable to unmarshal gamestate %s\n", err.Error())
	}


	//fmt.Printf("state %v\n", state)

	/*
		// write type...
		w, err := c.Writer(ctx, typ)
		if err != nil {
			return err
		}
	*/

	//w.Write(buf.Bytes())
	/*
		c.Writer(ctx, )
		_, err = io.Copy(w, r)
		if err != nil {
			return fmt.Errorf("failed to io.Copy: %w", err)
		}
	*/

	//err = w.Close()
	return err
}
