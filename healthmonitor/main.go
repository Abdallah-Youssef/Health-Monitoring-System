package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/robfig/cron/v3"
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

//flush the msgs to HDFS at midnight to save it in the aassociated log file
func flush(batch []HealthMessage, msg_counter int) {
	fmt.Printf("flushed %d", msg_counter)
}
func main() {
	//corn is used to schedule a function to run ...@midnight to flush the batch and add it to it's day

	p := make([]byte, 2048)
	msg_batch := make([]HealthMessage, 1024)
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
	cron_job := cron.New()
	cron_job.AddFunc("@midnight", func() { flush(msg_batch, int(msg_counter)) })
	cron_job.Start()

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
		//add the msg to the batch
		msg_batch[msg_counter] = recieved_msg
		print_msg_content(msg_batch[msg_counter])
		//if there is no error increase msg counter to keep track of rercieved msgs
		msg_counter++
		fmt.Printf("%d", msg_counter)
		if msg_counter == 1024 {
			//a batch is formed

		}
	}
}
