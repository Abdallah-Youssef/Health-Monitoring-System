package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/colinmarc/hdfs"
	"github.com/robfig/cron"
)

type Mem struct {
	Total float32
	Free  float32
}

type HealthMessage struct {
	ServiceName string
	TimeStamp   int64
	CPU         float32
	RAM         Mem
	Disk        Mem
}

var msg_batch [1024]HealthMessage
var msg_counter int = 0
var hdfsClient *hdfs.Client
var logFile *hdfs.FileWriter
var timeSinceLastFlush = time.Now()
var recieveTimes [1024]int64

// returns "/dd_mm_yyyy.log" of current day
func getDateString() string {
	t := time.Now()
	return fmt.Sprintf("/%v_%v_%v.log", t.Day(), int(t.Month()), t.Year())
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

func getRecieveTimesStats() {
	var sum int64
	var sd, mean float64
	now := time.Now().UnixNano()
	for _, t := range recieveTimes {
		sum += now - t
	}

	mean = float64(sum) / float64(msg_counter)

	for _, t := range recieveTimes {
		sd += math.Pow(float64(now-t)-mean, 2)
	}

	sd = math.Sqrt(sd / float64(msg_counter))
	fmt.Printf("Average wait time %v, std %v\n", mean, sd)
}

// flush whatever is in the msg_batch
func flush() {
	getRecieveTimesStats()
	fmt.Printf("Time taken to recieve 1024 msgs: %v\n", time.Since(timeSinceLastFlush))
	timeSinceLastFlush = time.Now()

	writeTime := time.Now()
	for i := 0; i < msg_counter; i++ {
		b, err := json.Marshal(msg_batch[i])
		if err != nil {
			fmt.Print("Failed to encode msg_batch\n")
		} else {
			io.WriteString(logFile, string(b))
		}
	}
	logFile.Flush()

	fmt.Printf("Time taken to flush 1024: %v\n", time.Since(writeTime))
	msg_counter = 0
}

func recieve_msg(p []byte, num_recieved_bytes int, m *sync.Mutex) {
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

	//add the msg to the batch
	m.Lock()
	msg_batch[msg_counter] = recieved_msg
	recieveTimes[msg_counter] = time.Now().UnixNano()
	msg_counter++
	fmt.Println(msg_counter)
	if msg_counter == 1024 {
		//a batch is formed
		flush()
	}
	m.Unlock()
}

func startUDPServer() *net.UDPConn {
	addr := net.UDPAddr{
		Port: 3500,
		IP:   net.ParseIP("127.0.0.1"),
	}

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Failed to start udp server: %v\n", err)
		os.Exit(1)
	}

	return ser
}

func connecToHDFS() {
	client, err := hdfs.New(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to connect to hdfs: %v\n", err)
		os.Exit(0)
	}
	hdfsClient = client
}

func openTodaysLogFile() {
	// Open new file for the day
	fileName := getDateString()
	fmt.Println(fileName)
	file, err := hdfsClient.Append(fileName)
	if err != nil {
		fmt.Printf("Failed to append to %v :\n, trying to Create", fileName)

		file, err = hdfsClient.Create(fileName)
		if err != nil {
			fmt.Printf("Failed to create %v :\n %v", fileName, err)
			return
		}

	}

	logFile = file
}

func main() {
	// Close and Flush logFile in the event of KeyBoardInterrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logFile.Close()
		os.Exit(2)
	}()

	//corn is used to schedule a function to run ...@midnight to flush the batch and add it to it's day
	ser := startUDPServer()
	connecToHDFS()
	openTodaysLogFile()
	defer logFile.Close()

	cron_job := cron.New()
	cron_job.AddFunc("@midnight", func() {
		// Flush any in-memory messages
		flush()

		// Close the previous file
		if logFile != nil {
			logFile.Close()
		}

		openTodaysLogFile()
	})
	cron_job.Start()
	var mu sync.Mutex
	for {
		//recieve msg in a byte array p
		p := make([]byte, 6*1024)
		num_recieved_bytes, _, err := ser.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		go recieve_msg(p, num_recieved_bytes, &mu)
	}
}
