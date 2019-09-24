package main

import (
	"bufio"
	"flag"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	port := flag.String("listen", ":25565", "The port / IP combo you want to listen on")
	flag.Parse()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Fatal error config file: %s", err)
	}
	viper.WatchConfig()

	lis, err := net.Listen("tcp", *port)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	logrus.Infof("Starting Listener")
	for {
		con, err := lis.Accept()
		if err != nil {
			logrus.Printf("Unable to accept a connection: %s", err)
			continue
		}
		go HandleConnection(con)
	}
}

func HandleConnection(con net.Conn) {
	defer con.Close()
	r := bufio.NewReader(con)

	// Read incoming packet
	packetId, data, packet := ReadPacket(r)

	if packetId == 0x00 {
		HandleHandshake(con, data, packet)
	}

}