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

var serviceName string
var freq int = 1000 // milliseconds
var healthMonitorAddr = "127.0.0.1:3500"

// This struct can be used for both the "RAM" and "DISK" fields
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
		ServiceName: serviceName,
		TimeStamp:   time.Now().UnixNano(),
		CPU:         getCPUUsage(),
		RAM:         getRamUsage(),
		Disk:        getDiskUsage(),
	}
}

func padString(str string) string {
	if len(str) > 200 {
		fmt.Printf("WARNING PACKET SIZE (%v) > 200\n INCREASE BUFFER SIZE\n EXITING", len(str))
		os.Exit(1)
	}

	ret := str
	count := 0
	for len(ret) < 200 {
		ret += " "
		count++
	}

	return ret
}


//  Args: freq in ms, health
func main() {
	serviceName, err := os.Hostname()
	if err != nil {
		fmt.Printf("Failed to get hostname")
		return
	}else {
		fmt.Printf("Host name: %v\n", serviceName)
	}

	// Parse frequency
	if len(os.Args) >= 2 {
		x, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Cannot convert \"%s\" to milliseconds\n", os.Args[1])
			os.Exit(1)
		}

		freq = x
	}


	
	healthMonitorAddr, err := net.ResolveUDPAddr("udp", "healthmonitor")
	if err != nil {
		fmt.Printf("Failed to get resolve healthmonitor's ip")
		return
	}else {
		fmt.Printf("Found health montior at %v\n", healthMonitorAddr)
	}
	
	conn, err := net.DialUDP("udp", nil, healthMonitorAddr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	defer conn.Close()

	for {
		healthMessage := getHealthJSON()
		bytes, _ := json.Marshal(healthMessage)

		if len(bytes) != 0 {
			str := padString(string(bytes))
			fmt.Fprintf(conn, str)
		}
	}

}
