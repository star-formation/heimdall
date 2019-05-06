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
package main

import (
  "os"

  "github.com/ethereum/go-ethereum/log"
  "github.com/urfave/cli"
  
  "github.com/star-formation/heimdall/websocket"
)

func init() {
	log.Root().SetHandler(log.MultiHandler(
		log.StreamHandler(os.Stderr, log.TerminalFormat(true)),
		log.LvlFilterHandler(
			log.LvlDebug,
			log.Must.FileHandler("heimdall_errors.json", log.JSONFormat()))))
}

func main() {
  app := cli.NewApp()
  app.Name = "heimdall"
  app.Version = "0.0.1"
  app.Usage = "The guardian of the gods.  Heimdall will blow a horn, called the Gjallarhorn, if Asgard is in danger."
  app.Action = func(c *cli.Context) error {
	log.Info("heimdall starting..")
	websocket.Debug()
    return nil
  }
  
  err := app.Run(os.Args)
  if err != nil {
    log.Error("app.Run:", "err", err)
  }
}