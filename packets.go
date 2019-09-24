package main

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/sirupsen/logrus"
)

// ReadPacket returns the packetId, the unread data, and the full packet
func ReadPacket(r *bufio.Reader) (uint, []byte, []byte) {
	length, err := binary.ReadUvarint(r)
	if err != nil {
		if err != io.EOF {
			logrus.Printf("An error occurred when reading packet length: %d", err)
		}
		return 0, nil, nil
	}
	packet := make([]byte, length)
	bytesRead, _ := io.ReadFull(r, packet)

	if int(length) != bytesRead {
		fullPacket := append(MakeVarInt(int(length)), packet...)
		logrus.Printf("Received unknown packet, proceeding as legacy packet 0x%x", length)
		return uint(length), packet, fullPacket
	}

	// Read packet id
	packetId, bytesRead := ReadVarInt(packet)
	if bytesRead <= 0 {
		logrus.Printf("An error occurred when reading packet id of packet: %d", bytesRead)
		return 0, nil, nil
	}
	i := bytesRead

	if length == 0 {
		return uint(packetId), []byte{}, append(MakeVarInt(int(length)), packet...)
	} else {
		return uint(packetId), packet[i:], append(MakeVarInt(int(length)), packet...)
	}
}

func MakePacket(packetId int, data []byte) []byte {
	packet := append(MakeVarInt(packetId), data...)
	return append(MakeVarInt(len(packet)), packet...)
}

func ReadVarInt(data []byte) (int, int) {
	value, bytesRead := binary.Uvarint(data)
	if bytesRead <= 0 {
		logrus.Printf("An error occurred while reading VarInt: %d", bytesRead)
		return 0, bytesRead
	}
	return int(value), bytesRead
}

func MakeVarInt(value int) []byte {
	temp := make([]byte, 10)
	bytesWritten := binary.PutUvarint(temp, uint64(value))
	return temp[:bytesWritten]
}

func ReadString(data []byte) (string, int) {
	length, bytesRead := ReadVarInt(data)
	if bytesRead <= 0 {
		logrus.Printf("An error occurred while reading string: %d", bytesRead)
		return "", bytesRead
	}
	return string(data[bytesRead : bytesRead+length]), bytesRead + length
}

func MakeString(str string) []byte {
	data := []byte(str)
	return append(MakeVarInt(len(data)), data...)
}
