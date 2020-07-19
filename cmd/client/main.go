package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kpfaulkner/websocketproxy/common"

	"nhooyr.io/websocket"
)

func readSocket(ctx context.Context, c *websocket.Conn) error {

	// read forevers...
	for {
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

		fmt.Printf("returned state %v\n", state)
	}

	return nil
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://10.0.0.100:8081", &websocket.DialOptions{})
	if err != nil {
		fmt.Printf("error dialing %s\n", err.Error())
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	// go and read the socket.
	go readSocket(ctx, c)

	state := common.PlayerState{PlayerID: 1, ChangeType: 2, X: 100, Y: 200, Dir: 7}
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("error Marshalling %s\n", err.Error())
	}


	start := time.Now()
	for i := 0; i < 10000; i++ {
		dataStr := time.Now().String()
		data = []byte(dataStr)
		err = c.Write(ctx, websocket.MessageBinary, data)
		time.Sleep(100 * time.Millisecond)
	}
	end := time.Now()

	fmt.Printf("took %dms\n", end.Sub(start).Milliseconds())

	fmt.Printf("going to sleep\n")
	time.Sleep(1000 * time.Second)
	c.Close(websocket.StatusNormalClosure, "")
}
