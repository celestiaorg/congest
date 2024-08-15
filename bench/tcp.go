package bench

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Message struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}

type server struct {
}

func (s *server) Run(port string, bufferSize int) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("TCP server started on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.SendRoutine(conn)
	}
}

func (s *server) Stop() {

}

func (s *server) SendRoutine(conn net.Conn) {
	for {
		// create a random 64KB message
		data := make([]byte, 1024*64)
		_, err := rand.Read(data)
		if err != nil {
			fmt.Println("Error generating random data:", err)
			panic(err)
		}

		// Send response with current timestamp
		msg := Message{
			Timestamp: time.Now(),
			Data:      data,
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshalling message:", err)
			panic(err)
		}

		// Get the length of the message
		msgLength := uint32(len(msgBytes))

		// Create a buffer to hold the length prefix and the message
		buf := make([]byte, 4+len(msgBytes))

		// Write the length prefix to the buffer
		binary.BigEndian.PutUint32(buf[:4], msgLength)

		// Write the message to the buffer
		copy(buf[4:], msgBytes)

		conn.Write(buf)
	}
}

func (s *server) ReadRoutine(conn net.Conn) {
	for {
		// Read the length prefix (4 bytes)
		lengthBuf := make([]byte, 4)
		_, err := conn.Read(lengthBuf)
		if err != nil {
			fmt.Println("Error reading length prefix:", err)
			return
		}

		// Unpack the length prefix to get the message length
		msgLength := binary.BigEndian.Uint32(lengthBuf)

		// Read the message of the specified length
		msgBytes := make([]byte, msgLength)
		_, err = conn.Read(msgBytes)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Deserialize the message from JSON
		var msg Message
		err = json.Unmarshal(msgBytes, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}
	}
}
