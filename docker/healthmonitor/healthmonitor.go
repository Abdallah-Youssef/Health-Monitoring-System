package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/colinmarc/hdfs"
	"github.com/robfig/cron"
)

var msg_batch [1024]HealthMessage
var msg_counter int = 0

func StartHealthMonitor(hdfsAddr string) {
	udpServer := startUDPServer()
	hdfsClient := connecToHDFS(hdfsAddr)
	logFile := openTodaysLogFile(hdfsClient)

	cron_job := cron.New()
	cron_job.AddFunc("@midnight", func() {
		// Flush any in-memory messages
		flush(logFile)

		// Close the previous file
		if logFile != nil {
			logFile.Close()
		}

		logFile = openTodaysLogFile(hdfsClient)
	})
	cron_job.Start()

	for {
		//recieve msg in a byte array p
		p := make([]byte, 200)
		num_recieved_bytes, _, err := udpServer.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		fmt.Println(num_recieved_bytes)

		recieve_msg(p, logFile)
	}
}

func recieve_msg(p []byte, logFile *hdfs.FileWriter) {
	//convert byte array into json object
	var recieved_msg HealthMessage
	var error_unmarshal = json.Unmarshal(p, &recieved_msg)
	if error_unmarshal != nil {
		fmt.Printf("error in parsing json  %v\n", error_unmarshal)
		fmt.Println("\"", string(p), "\"")
		os.Exit(1)
		return
	}

	//add the msg to the batch
	msg_batch[msg_counter] = recieved_msg
	// recieveTimes[msg_counter] = time.Now().UnixNano()
	msg_counter++
	fmt.Print(msg_counter, "-", recieved_msg.ServiceName)
	if msg_counter == 1024 {
		//a batch is formed
		// flush(logFile)
	}
}

func startUDPServer() *net.UDPConn {
	addr := net.UDPAddr{
		Port: 3500,
		IP:   net.ParseIP("0.0.0.0"),
	}

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Failed to start udp server: %v\n", err)
		os.Exit(1)
	}

	return ser
}

func connecToHDFS(hdfsAddr string) *hdfs.Client {
	fmt.Printf("Connecting to hdfs node at %v:", hdfsAddr)
	client, err := hdfs.New(hdfsAddr)
	if err != nil {
		fmt.Printf("Failed to connect to hdfs: %v\n", err)
		os.Exit(0)
	}

	fmt.Printf("Connected successfully\n")
	return client
}

// returns "/dd_mm_yyyy.log" of current day
func getDateString() string {
	t := time.Now()
	return fmt.Sprintf("/%v_%v_%v.log", t.Day(), int(t.Month()), t.Year())
}

func openTodaysLogFile(hdfsClient *hdfs.Client) *hdfs.FileWriter {
	// Open new file for the day
	fileName := getDateString()
	fmt.Println(fileName)
	file, err := hdfsClient.Append(fileName)
	if err != nil {
		fmt.Printf("Failed to append to \"%v\": %v\n\n, trying to create\n", fileName, err)

		file, err = hdfsClient.Create(fileName)
		if err != nil {
			fmt.Printf("Failed to create %v :\n %v", fileName, err)
			os.Exit(0)
		}

	}

	// register an interrupt handler that makes sure that the file content is written before terminating
	registerInterruptHandler(file)
	return file
}

func registerInterruptHandler(logFile *hdfs.FileWriter) {
	// Close and Flush logFile in the event of KeyBoardInterrupt
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logFile.Close()
		os.Exit(2)
	}()
}

// flush whatever is in the msg_batch
// The commented code is from the time measurements for the report
func flush(logFile *hdfs.FileWriter) {
	//getRecieveTimesStats()
	//fmt.Printf("Time taken to recieve 1024 msgs: %v\n", time.Since(timeSinceLastFlush))
	//timeSinceLastFlush = time.Now()

	//writeTime := time.Now()
	for i := 0; i < msg_counter; i++ {
		b, err := json.Marshal(msg_batch[i])
		if err != nil {
			fmt.Print("Failed to encode msg_batch\n")
		} else {
			io.WriteString(logFile, string(b))
		}
	}
	logFile.Flush()

	//fmt.Printf("Time taken to flush 1024: %v\n", time.Since(writeTime))
	msg_counter = 0
}

// var timeSinceLastFlush = time.Now()
// var recieveTimes [1024]int64
// func getRecieveTimesStats() {
// 	var sum int64
// 	var sd, mean float64
// 	now := time.Now().UnixNano()
// 	for _, t := range recieveTimes {
// 		sum += now - t
// 	}
// 	mean = float64(sum) / float64(msg_counter)
// 	for _, t := range recieveTimes {
// 		sd += math.Pow(float64(now-t)-mean, 2)
// 	}
// 	sd = math.Sqrt(sd / float64(msg_counter))
// 	fmt.Printf("Average wait time %v, std %v\n", mean, sd)
// }

func print_msg_content(msg HealthMessage) {
	fmt.Printf("service name %s \n", msg.ServiceName)
	fmt.Printf("Timestamp %d \n", msg.TimeStamp)
	fmt.Printf("CPU %f \n", msg.CPU)
	fmt.Printf("Ram total %f \n", msg.RAM.Total)
	fmt.Printf("Ram free %f \n", msg.RAM.Free)
	fmt.Printf("Disk total %f \n", msg.Disk.Total)
	fmt.Printf("Disk free %f \n", msg.Disk.Free)
}
