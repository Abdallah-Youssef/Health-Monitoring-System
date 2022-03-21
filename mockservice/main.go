package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/mackerelio/go-osstat/cpu" // for cpu usage
	"github.com/pbnjay/memory"            // for ram usage
	"github.com/ricochet2200/go-disk-usage/du"
)

var serviceName string = "Default Service Name"
var freq int // milliseconds

// This struct can be used for both the "RAM" and "DISK" fields
type Mem struct {
	Total float32
	Free  float32
}

type HealthMessage struct {
	serviceName string
	TimeStamp   int64
	CPU         float32
	RAM         Mem
	Disk        Mem
}

func getRamUsage() Mem {
	var total float32 = float32(memory.TotalMemory()) / 1024 / 1024 / 1024
	var free float32 = float32(memory.FreeMemory()) / 1024 / 1024 / 1024
	return Mem{
		Total: total,
		Free:  free,
	}
}

func getCPUUsage() float32 {
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return -1
	}
	time.Sleep(time.Duration(freq) * time.Millisecond)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return -1
	}

	total := float32(after.Total - before.Total)
	return float32(after.System-before.System) / total
}

func getDiskUsage() Mem {
	diskUsage := du.NewDiskUsage("/")

	total := float32(diskUsage.Size()) / 1024 / 1024 / 1024
	free := float32(diskUsage.Available()) / 1024 / 1024 / 1024
	return Mem{
		Total: total,
		Free:  free,
	}
}

func getHealthJSON() HealthMessage {
	return HealthMessage{
		serviceName: serviceName,
		TimeStamp:   time.Now().UnixNano(),
		CPU:         getCPUUsage(),
		RAM:         getRamUsage(),
		Disk:        getDiskUsage(),
	}
}

func main() {
	// Parse frequency
	if len(os.Args) >= 2 {
		x, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Cannot convert \"%s\" to milliseconds\n", os.Args[1])
			os.Exit(1)
		}

		freq = x
	}

	if len(os.Args) >= 3 {
		serviceName = os.Args[2]
	}

	conn, err := net.Dial("udp", "127.0.0.1:3500")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	for {
		healthMessage := getHealthJSON()
		str, _ := json.Marshal(healthMessage)

		fmt.Println(string(str))
		fmt.Fprintf(conn, string(str))
	}

}
