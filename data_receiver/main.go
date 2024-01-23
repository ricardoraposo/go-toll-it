package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	t "github.com/ricardoraposo/toll-calculator/types"
)

var u = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type DataReceiver struct {
	msg  chan t.OBUData
	conn *websocket.Conn
}

func main() {
	receiver := NewDataReceiver()

	http.HandleFunc("/ws", receiver.handleWS)

	http.ListenAndServe(":30000", nil)
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msg: make(chan t.OBUData, 128),
	}
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected !!!")
	for {
		var data t.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("Read error: ", err)
			continue
		}
		fmt.Printf("Received OBU data from [%d] :: <Lat %.2f, Long %.2f> \n", data.OBUID, data.Lat, data.Long)
		dr.msg <- data
	}
}
