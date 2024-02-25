package ws

import (
	"errors"
	"log"

	"wQueue/db"

	"github.com/gorilla/websocket"
)

type WConn struct {
	Conn *websocket.Conn
	Qid  int
}

func MessageWs(cons []WConn, q *db.Queue, s string) {
	for _, c := range cons {
		if c.Qid == q.Id {
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte("M"+s))
			if err != nil {
				if !errors.Is(err, websocket.ErrCloseSent) {
					log.Print(err)
				}
			}
		}
	}
}

func ReloadWs(cons []WConn, q *db.Queue) {
	for _, c := range cons {
		if c.Qid == q.Id {
			err := c.Conn.WriteMessage(websocket.TextMessage, []byte("R"))
			if err != nil {
				if !errors.Is(err, websocket.ErrCloseSent) {
					log.Print(err)
				}
			}
		}
	}
}

func RemoveConnection(cons []WConn, conn *websocket.Conn) []WConn {
	for i, c := range cons {
		if c.Conn == conn {
			cons = append(cons[:i], cons[i+1:]...)
			return cons
		}
	}
	return cons
}
