package main

import (
	"flag"
	"strconv"

	"github.com/makizm/ssmartthings/lan-device-handler/routing"
)

func main() {
	srvPort := flag.Int("port", 3000, "server port number to listen on")
	srvMode := flag.String("mode", "debug", "Server running mode: debug or prod")
	flag.Parse()

	server := engine.APIServer(*srvMode)
	server.Run(":" + strconv.Itoa(*srvPort))
}
