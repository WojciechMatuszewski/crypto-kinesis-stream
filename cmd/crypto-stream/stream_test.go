package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

var upgrader = websocket.Upgrader{}

func TestStream_Init(t *testing.T) {
	t.Run("initialize connection success", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf(err.Error())
			}

			_, m, err := c.ReadMessage()
			if err != nil {
				t.Fatalf(err.Error())
			}

			assert.Equal(t, "{\"op\":\"unconfirmed_sub\"}", string(m))
			wg.Done()

		})
		server := httptest.NewServer(h)
		defer server.Close()

		u := url.URL{
			Scheme: "ws",
			Host:   strings.Replace(server.URL, "http://", "", -1),
		}

		stream := NewStream(func(s *Stream) {
			s.addr = u.String()
		})

		err := stream.Init()
		wg.Wait()

		assert.NoError(t, err)
	})

	t.Run("wrong addr", func(t *testing.T) {
		stream := NewStream(func(s *Stream) {
			s.addr = "123"
		})

		err := stream.Init()
		assert.Error(t, err)
	})
}

func TestStream_Listener(t *testing.T) {
	t.Run("not initialized", func(t *testing.T) {
		s := NewStream()
		_, err := s.GetListener()

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotInitialized))
	})

	t.Run("single transaction", func(t *testing.T) {
		var wg sync.WaitGroup
		trIn, trB := loadTestData(t, "single.json")

		wg.Add(1)
		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf(err.Error())
			}

			err = c.WriteMessage(websocket.TextMessage, trB)
			if err != nil {
				t.Fatalf(err.Error())
			}
			wg.Done()
		})

		server := httptest.NewServer(h)
		defer server.Close()

		u := url.URL{
			Scheme: "ws",
			Host:   strings.Replace(server.URL, "http://", "", -1),
		}

		stream := NewStream(func(s *Stream) {
			s.addr = u.String()
		})

		_ = stream.Init()
		listener, _ := stream.GetListener()

		tr, err := listener()
		wg.Wait()

		assert.NoError(t, err)
		assert.Equal(t, trIn, tr)
	})

	t.Run("multiple transactions", func(t *testing.T) {
		trIn, trB := loadTestData(t, "single.json")

		h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer c.Close()

			err = c.WriteMessage(websocket.TextMessage, trB)
			if err != nil {
				t.Fatalf(err.Error())
			}

			err = c.WriteMessage(websocket.TextMessage, trB)
			if err != nil {
				t.Fatalf(err.Error())
			}
		})

		server := httptest.NewServer(h)
		defer server.Close()

		u := url.URL{
			Scheme: "ws",
			Host:   strings.Replace(server.URL, "http://", "", -1),
		}

		stream := NewStream(func(s *Stream) {
			s.addr = u.String()
		})
		_ = stream.Init()

		var transactions []Transaction

		listener, _ := stream.GetListener()
		// probably very suboptimal
		for i := 0; i < 2; {
			trOut, err := listener()
			assert.NoError(t, err)

			transactions = append(transactions, trOut)
			i++
		}

		assert.Len(t, transactions, 2)
		assert.Equal(t, transactions, []Transaction{trIn, trIn})
	})
}

func loadTestData(t *testing.T, fileName string) (Transaction, []byte) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf(err.Error())
	}

	f, err := os.Open(path.Join(wd, "/testdata/", fileName))
	if err != nil {
		t.Fatalf(err.Error())
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tr, err := transactionFromBytes(b)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return tr, b
}
