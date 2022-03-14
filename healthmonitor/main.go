package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Mem struct {
	Total float32
	Free  float32
}

type HealthMessage struct {
	ServiceName string `json:serviceName`
	TimeStamp   int64
	CPU         float32
	RAM         Mem
	Disk        Mem
}

func print_msg_content(msg HealthMessage) {
	fmt.Printf("service name %s \n", msg.ServiceName)
	fmt.Printf("Timestamp %d \n", msg.TimeStamp)
	fmt.Printf("CPU %f \n", msg.CPU)
	fmt.Printf("Ram total %f \n", msg.RAM.Total)
	fmt.Printf("Ram free %f \n", msg.RAM.Free)
	fmt.Printf("Disk total %f \n", msg.Disk.Total)
	fmt.Printf("Disk free %f \n", msg.Disk.Free)
}

func main() {
	p := make([]byte, 2048)
	var msg_counter uint16 = 0
	addr := net.UDPAddr{
		Port: 3500,
		IP:   net.ParseIP("127.0.0.1"),
	}

	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	for {
		//recieve msg in a byte array p
		num_recieved_bytes, remoteaddr, err := ser.ReadFromUDP(p)
		if err != nil {

			fmt.Printf("Some error  %v", err)
			continue
		}
		// slice the byte array into the size of recieved bytes and ignore the rest
		p = p[:num_recieved_bytes]
		//convert byte array into json object
		var recieved_msg HealthMessage
		var error_unmarshal = json.Unmarshal(p, &recieved_msg)
		if error_unmarshal != nil {
			fmt.Printf("error in parsing json  %v", error_unmarshal)
			continue
		}
		fmt.Printf("Read a message from %v  \n", remoteaddr)
		print_msg_content(recieved_msg)
		//if there is no error increase msg counter to keep track of rercieved msgs
		msg_counter++

	}
}
