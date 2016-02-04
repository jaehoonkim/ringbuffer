package main

import (
	//"bytes"
	//"encoding/binary"

	"log"
	"os"
	"os/signal"

	rb "github.com/jaehoonkim/ringbuffer"
)

func ConsumerFunc(rb *rb.RingBuffer, index, size int) {
	for i := 0; i < size; i++ {
		log.Println("[", index, "]: ", "consumer : ", i, ":", rb.Read())
		rb.Read()
	}
}

func ProducerFunc(rb *rb.RingBuffer, index, size int) {
	for i := 0; i < size; i++ {
		rb.Write(i)
	}

}

func main() {
	rb, err := rb.NewRingBuffer(4)
	if err != nil {
		log.Println(err)
		return
	}

	go ProducerFunc(rb, 1, 8)
	go ConsumerFunc(rb, 1, 8)

	/*go ProducerFunc(rb, 1, 8)
	go ConsumerFunc(rb, 2, 7)*/
	/*
		rb.Write(1)
		rb.Write(2)
		rb.Write(3)
		rb.Write(4)
		rb.Write(5)
		rb.Write(6)
		rb.Write(7)
		rb.Write(8)
		rb.Write(9)

		fmt.Println("read : ", rb.Read())
		fmt.Println("read : ", rb.Read())
		fmt.Println("read : ", rb.Read())
	*/
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
