package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrNotInitialized = errors.New("stream not initialized")
)

var (
	defaultAddr    = "wss://ws.blockchain.info/inv"
	defaultTimeout = time.Second * 10
)

// Stream represents stream of crypto transactions
type Stream struct {
	timeout time.Duration
	addr    string
	conn    *websocket.Conn
}

// Option represents function which modifies the Stream
type Option func(*Stream)

// NewStream returns stream of crypto transactions
func NewStream(options ...Option) Stream {
	s := Stream{
		timeout: defaultTimeout,
		addr:    defaultAddr,
	}

	for _, opt := range options {
		opt(&s)
	}

	return s
}

// Init initializes the stream
func (s *Stream) Init() error {
	c, _, err := websocket.DefaultDialer.Dial(s.addr, nil)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"op":"unconfirmed_sub"}`))
	if err != nil {
		return err
	}

	s.conn = c
	return nil
}

// GetListener returns a listener which returns transactions from the stream
func (s Stream) GetListener() (func() (Transaction, error), error) {
	if s.conn == nil {
		return func() (Transaction, error) {
			return Transaction{}, nil
		}, ErrNotInitialized
	}

	return func() (Transaction, error) {
		var tr Transaction
		_, m, err := s.conn.ReadMessage()
		if err != nil {
			return tr, err
		}

		err = json.Unmarshal(m, &tr)
		if err != nil {
			return tr, err
		}

		return tr, nil
	}, nil
}

// Stop stops the listening to a stream
func (s Stream) Stop() error {
	if s.conn == nil {
		return nil
	}

	err := s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	err = s.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
