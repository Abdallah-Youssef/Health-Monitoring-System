package main

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
