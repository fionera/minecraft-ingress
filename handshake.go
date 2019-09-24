package main

import (
	"io"
	"net"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func HandleHandshake(con net.Conn, data []byte, packet []byte) {
	var i int
	// Handshake
	logrus.Printf("<%s> Received handshake", con.RemoteAddr().String())

	// Protocol version
	_, bytesRead := ReadVarInt(data[i:])
	if bytesRead <= 0 {
		logrus.Printf("<%s> An error occurred when reading protocol version of handshake packet: %d", con.RemoteAddr().String(), bytesRead)
		return
	}
	i += bytesRead

	// Address
	var address string
	address, bytesRead = ReadString(data[i:])
	if bytesRead <= 0 {
		logrus.Printf("<%s> An error occurred when reading server address of handshake packet: %d", con.RemoteAddr().String(), bytesRead)
		return
	}
	i += bytesRead

	// Port
	//port := binary.BigEndian.Uint16(data[i:i + 2])
	i += 2

	server := viper.GetStringMapString("server")

	if backendAddress, ok := server[address]; ok {
		if !strings.Contains(backendAddress, ":") {
			backendAddress += ":25565"
		}
		backendConnection, err := net.Dial("tcp", backendAddress)
		if err != nil {
			logrus.Printf("<%s> Error while connecting to the backend server: %s", con.RemoteAddr().String(), err)
			return
		}
		defer backendConnection.Close()

		if err != nil {
			logrus.Printf("<%s> Error while connecting to the backend server: %s", con.RemoteAddr().String(), err)
			con.Write(MakePacket(0x00, MakeString("\"Could not connect to the backend server. Please notify the server administrator.\"")))
			return
		}

		_, err = backendConnection.Write(packet)
		if err != nil {
			logrus.Printf("<%s> Error while relaying data to the backend server: %s", con.RemoteAddr().String(), err)
			con.Write(MakePacket(0x00, MakeString("\"Could not relay data to the backend server. Please notify the server administrator.\"")))
			return
		}

		go io.Copy(backendConnection, con)
		io.Copy(con, backendConnection)
	}


}
