package main

import (
	"time"
	"wb-tech-l0/internal/delivery/kafka"
)

func main() {
	go func() {
		time.Sleep(3 * time.Second)
		run_fake_data_producer()
	}()
	kafka.Init_consumer()
}
