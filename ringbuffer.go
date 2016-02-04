package ringbuffer

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

type RingBuffer struct {
	buffer   []int
	head     int64
	tail     int64
	size     int64 //buf의 전체 사이즈
	count    int64 //현재 버퍼에 차 있는 데이터 개수
	mask     int64
	producer *sync.Cond
	consumer *sync.Cond
}

func checkPowerOfTwo(size int64) bool {
	return ((size != 0) && (size&(size-1)) == 0)
}

func NewRingBuffer(bufSize int64) (*RingBuffer, error) {
	if ok := checkPowerOfTwo(bufSize); ok == false {
		return nil, errors.New("it is not power of two")
	}

	return &RingBuffer{
		buffer:   make([]int, bufSize),
		head:     0,
		tail:     0,
		size:     bufSize,
		count:    0,
		mask:     bufSize - 1,
		producer: sync.NewCond(&sync.Mutex{}),
		consumer: sync.NewCond(&sync.Mutex{}),
	}, nil
}

func (r *RingBuffer) GetBuffer() []int {
	return r.buffer
}

func (r *RingBuffer) Write(buf int) {
	r.producer.L.Lock()
	if r.IsFull() == true {
		//log.Println("producer wait..start")
		r.producer.Wait()
		log.Println("[IsFull]", "head: ", r.head, "tail: ", r.tail, "count: ", r.count, "size: ", r.size)
		//log.Println("producer wait..end")
	}

	r.buffer[r.head] = buf
	r.head = (r.head + 1) & r.mask

	r.count = atomic.AddInt64(&r.count, 1)

	r.consumer.Broadcast()
	r.producer.L.Unlock()

	//log.Println("[write]: ", buf, "[count]: ", r.count, "[head]: ", r.head, "[tail]: ", r.tail)
}

func (r *RingBuffer) Writes(buf []int) {

}

func (r *RingBuffer) Read() int {
	r.consumer.L.Lock()
	if r.IsEmpty() == true {
		log.Println("consumer wait...start")
		r.consumer.Wait()
		log.Println("consumer wait...end")
		//log.Println("[IsEmpty] ", "head: ", r.head, "tail: ", r.tail, "count: ", r.count, "size: ", r.size)
	}

	var buf int
	buf = r.buffer[r.tail]
	r.tail = (r.tail + 1) & r.mask

	r.count = atomic.AddInt64(&r.count, -1)

	//log.Println("tail: ", r.tail)
	//log.Println("count: ", r.count)

	r.producer.Broadcast()
	r.consumer.L.Unlock()

	log.Println("[read]: ", buf, "[count]: ", r.count, "[head]: ", r.head, "[tail]: ", r.tail)

	return buf
}

func (r *RingBuffer) Reads(size int) []int {
	var buf []int
	r.consumer.L.Lock()

	r.producer.L.Unlock()

	return buf

}

func (r *RingBuffer) IsFull() bool {
	//log.Println("head: ", r.head, "tail: ", r.tail, "count: ", r.count, "size: ", r.size)
	return (r.head == r.tail) && (r.count == r.size)
}

func (r *RingBuffer) IsEmpty() bool {
	//log.Println("[IsEmpty] ", "head: ", r.head, "tail: ", r.tail, "count: ", r.count, "size: ", r.size)
	return ((r.tail == r.head) && (r.count == 0))
}

/*
	입력하고자 하는 버퍼의 길이가 현재 버퍼상태에서 입력이 가능한지 체크한다.
	size : 입력될 버퍼의 길이
*/
func (r *RingBuffer) IsWriteAvailable() bool {
	index := r.head + 1
	index = index & r.mask

	var result bool
	if index < r.head && index > r.tail {
		result = false
	} else if index == r.tail { // full인 상태까지 입력됨
		result = true
	} else {
		result = true
	}

	return result
}

func (r *RingBuffer) IsReadAvailable(size int64) bool {
	index := r.tail + size
	index = index & r.mask

	var result bool
	if index >= r.head && index < r.tail {
		result = false
	} else {
		result = true
	}

	return result
}
