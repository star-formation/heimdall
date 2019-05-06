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

    "github.com/ethereum/go-ethereum/log"
    "github.com/gorilla/websocket"
)

// note: localhost:8081 as addr string works with raw TCP but not with HTTP
var host = ":8081"

var upgrader = websocket.Upgrader{} // use default options

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info("func handler: ")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade err:", "err", err)
		return
	}
	defer c.Close()
	for {
		log.Info("waiting on c.ReadMessage: ")
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error("read:", "err", err)
			break
		}
		log.Info("recv: ", "msg", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error("write err:", "err", err)
			break
		}
	}
}

func Debug() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		log.Error("http.ListenAndServe", "err", err)
	}
}
