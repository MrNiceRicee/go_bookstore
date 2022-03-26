package main

import (
	"fmt"
	"net"
	"net/http"
	"server/connection"
	"server/router"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	app := router.Routes()
	connection.CreateConnection()

	fmt.Printf("Local: http://%v:8000\nNetwork: http://%v:8000", "localhost", GetOutboundIP())
	http.ListenAndServe(":8000", app)
}
