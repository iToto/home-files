package coinroutespriceconsumer

import (
	"encoding/json"
	"fmt"
	"time"
	"yield-mvp/internal/wlog"

	"github.com/gorilla/websocket"
)

type Config struct {
	ResourcePath string
	Payload      interface{}
}

type Consumer struct {
	conf         Config
	lastResponse []byte
}

func NewConsumer(conf Config) (*Consumer, error) {

	// ensure mandatory config vars are set
	if conf.ResourcePath == "" {
		return nil, fmt.Errorf("missing mandatory config: Resource Path")
	}
	if conf.Payload == nil {
		return nil, fmt.Errorf("missing mandatory config: Payload")
	}

	return &Consumer{
		conf: conf,
	}, nil
}

func (c *Consumer) Start(
	conn *websocket.Conn,
	wl wlog.Logger,
) {
	go c.start(conn, wl)
}

func (c *Consumer) GetPrice() (float64, error) {
	r := RealPriceResponse{}
	err := json.Unmarshal(c.lastResponse, &r)
	if err != nil {
		return 0.0, err
	}

	return r.Price, nil
}

func (c *Consumer) GetLastResponse() string {
	return string(c.lastResponse)
}

func (c *Consumer) start(
	conn *websocket.Conn,
	wl wlog.Logger,
) {
	wl = wl.WithStr("websocketPath", c.conf.ResourcePath)
	done := make(chan struct{})

	go c.responseHandler(conn, done, wl)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			wl.Debugf("received channel closed, exiting...")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				wl.Infof("cannot close connection gracefully:", err)
				return
			}
			return
		case <-ticker.C:
			// make a request every second
			// wl.Debugf("sending payload: %+v", c.conf.Payload)
			// json, err := json.Marshal(c.conf.Payload)
			// if err != nil {
			// 	wl.Debugf("cannot marshall json", err)
			// 	return
			// }
			// wl.Debugf("sending message: %+v \n", string(json))
			err := conn.WriteJSON(c.conf.Payload)
			if err != nil {
				wl.Debugf("error sending message:", err)
				return
			}
		}
	}
}

func (c *Consumer) responseHandler(
	conn *websocket.Conn,
	done chan struct{},
	wl wlog.Logger,
) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			wl.Debugf("read error:", err)
			wl.Debugf("close done channel")
			return
		}
		// wl.Debugf("received: %s\n", message)

		c.lastResponse = message
	}
}
