package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

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
func flush() {
	fmt.Printf("flushed %d", msg_counter)
}
func recieve_msg(p []byte, wg *sync.WaitGroup, m *sync.Mutex, num_recieved_bytes int) {
	recieved_bytes := make([]byte, num_recieved_bytes)
	// slice the byte array into the size of recieved bytes and ignore the rest
	recieved_bytes = p[:num_recieved_bytes]
	//convert byte array into json object
	var recieved_msg HealthMessage
	var error_unmarshal = json.Unmarshal(recieved_bytes, &recieved_msg)
	if error_unmarshal != nil {
		fmt.Printf("error in parsing json  %v", error_unmarshal)
		return
	}
	//fmt.Printf("Read a message from %v  \n", remoteaddr)
	//add the msg to the batch
	m.Lock()
	msg_batch[msg_counter] = recieved_msg
	print_msg_content(msg_batch[msg_counter])
	//if there is no error increase msg counter to keep track of rercieved msgs
	msg_counter++

	fmt.Printf("%d", msg_counter)
	if msg_counter == 1024 {
		//a batch is formed
		msg_counter = 0

	}
	m.Unlock()
	wg.Done() //notify all the waiting

}

var msg_batch [1024]HealthMessage
var msg_counter uint16 = 0

func main() {
	//corn is used to schedule a function to run ...@midnight to flush the batch and add it to it's day

	p := make([]byte, 6*1024)

	var w sync.WaitGroup
	var m sync.Mutex
	addr := net.UDPAddr{
		Port: 3500,
		IP:   net.ParseIP("127.0.0.1"),
	}

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	fmt.Print(ser)
	cron_job := cron.New()
	cron_job.AddFunc("@midnight", func() { flush() })
	cron_job.Start()

	for count := 0; count <= 120; count++ {
		//recieve msg in a byte array p
		num_recieved_bytes, remoteaddr, err := ser.ReadFromUDP(p)

		if err != nil {

			fmt.Printf("Some error  %v", err)
			continue
		}
		fmt.Printf("Read a message from %v  \n", remoteaddr)
		w.Add(1)
		go recieve_msg(p, &w, &m, num_recieved_bytes)
		w.Wait()
	}
}
