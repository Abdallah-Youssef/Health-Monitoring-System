package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/colinmarc/hdfs"
)

const batch_buffer_size int = 10
var batch_messages_length int = 0
var batch_messages [batch_buffer_size]string


func StartHealthMonitor(hdfsAddr string) {
	udpServer := startUDPServer()
	hdfsClient := connecToHDFS(hdfsAddr)
	batchFile := openHDFSFile("/messages.csv", hdfsClient)
	speedFile := openHDFSFile("/new.csv", hdfsClient)

	for {
		//recieve msg in a byte array p
		p := make([]byte, 200)
		num_recieved_bytes, _, err := udpServer.ReadFromUDP(p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		fmt.Println(num_recieved_bytes)

		recieve_msg(p, batchFile, speedFile)
	}
}

func recieve_msg(p []byte, batchFile *hdfs.FileWriter, speedFile *hdfs.FileWriter) {
	//convert byte array into json object
	var recieved_msg HealthMessage
	var error_unmarshal = json.Unmarshal(p, &recieved_msg)
	if error_unmarshal != nil {
		fmt.Printf("error in parsing json  %v\n", error_unmarshal)
		fmt.Println("\"", string(p), "\"")
		os.Exit(1)
		return
	}

	print_msg_content(recieved_msg)

	//add the msg to the batch
	batch_messages[batch_messages_length] = toString(recieved_msg)

	batch_messages_length++

	if batch_messages_length >= batch_buffer_size {
		flush(batchFile, speedFile)
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


func openHDFSFile(fileName string, hdfsClient *hdfs.Client) *hdfs.FileWriter {
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

// flush whatever is in the batch
// The commented code is from the time measurements for the report
func flush(batchFile *hdfs.FileWriter, speedFile *hdfs.FileWriter) {
	//getRecieveTimesStats()
	//fmt.Printf("Time taken to recieve 1024 msgs: %v\n", time.Since(timeSinceLastFlush))
	//timeSinceLastFlush = time.Now()

	//writeTime := time.Now()
	for i := 0; i < batch_messages_length; i++ {
			io.WriteString(batchFile, batch_messages[i])
			io.WriteString(speedFile, batch_messages[i])
	}
	batchFile.Flush()
	speedFile.Flush()

	//fmt.Printf("Time taken to flush 1024: %v\n", time.Since(writeTime))
	batch_messages_length = 0
}


func toString(msg HealthMessage) string{
		var b bytes.Buffer
		b.WriteString(msg.ServiceName)
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v", msg.TimeStamp))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v",(msg.CPU)))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v",(msg.RAM.Total)))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v",(msg.RAM.Free)))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v",(msg.Disk.Total)))
		b.WriteString(",")
		b.WriteString(fmt.Sprintf("%v",(msg.Disk.Free)))
		b.WriteString("\n")

		fmt.Print(b.String())
		return b.String()
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
