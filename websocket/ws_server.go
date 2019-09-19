/*  Copyright 2019 The heimdall Authors

    This file is part of heimdall.

    heimdall is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    heimdall is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package websocket

import (
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/gorilla/websocket"
	"github.com/star-formation/tesseract"
)

// note: localhost:8081 as addr string works with raw TCP but not with HTTP
var host = ":8081"

var upgrader = websocket.Upgrader{} // use default options

func getHTTPHandler(protocolHandler func([]byte) ([]byte, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("func handler: ", "r", r)
		// TODO: secure origin check
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("upgrade err:", "err", err)
			return
		}
		defer c.Close()

		sps := websocket.Subprotocols(r)
		if len(sps) != 1 || sps[0] != "client0.argonavis.io" {
			log.Error("Unsupported WebSocket Subprotocol: ", "subs", sps)
			WriteControlClose(c, websocket.CloseProtocolError, "unsupported subprotocol")
			return
		}

		for {
			log.Info("waiting on c.ReadMessage: ")
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Error("read:", "err", err)
				break
			}
			log.Info("recv: ", "msg", msg)

			// setup MessageBus sub to engine loop
			ch := tesseract.S.MB.Subscribe()
			go func() {
				for {
					stateJSON := <-ch
					err = c.WriteMessage(mt, stateJSON)
					if err != nil {
						log.Error("write err:", "err", err)
						break
					}
				}
			}()

			/*
				resp, err := protocolHandler(msg)
				if err != nil {
					log.Error("protocolHandler", "err", err)
					WriteControlClose(c, websocket.CloseProtocolError, err.Error())
				}

				//log.Debug("c.WriteMessage", "resp", resp)
				err = c.WriteMessage(mt, resp)
				if err != nil {
					log.Error("write err:", "err", err)
					break
				}
			*/
		}
	}
}

func WriteControlClose(c *websocket.Conn, closeCode int, str string) error {
	msg := websocket.FormatCloseMessage(closeCode, str)
	return c.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second*5))
}

func Start(protocolHandler func([]byte) ([]byte, error)) {
	handler := getHTTPHandler(protocolHandler)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Error("http.ListenAndServe", "err", err)
	}
	log.Info("WebSocket Server Started", "host", host)
}
