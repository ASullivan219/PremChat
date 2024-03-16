package messagequeue

import (
	"fmt"
	"time"
)

type Iterator interface{
	Next()
	HasNext()
}


type Queue interface{
	AddMessage()
}

type MessageQueueIterator struct{
	Iterator
	Queue 			*MessageQueue
	currentIndex	int
	elementsLeft	int
	size			int
}

type Message struct{
	RawText		string
	Username	string
	Time		time.Time
}

func NewMessage( rawText string, username string) Message{
	currentTime := time.Now()
	return Message{ RawText: rawText, Username: username, Time: currentTime} 
}

type MessageQueue struct{
	Queue
	Messages 	[]Message
	First		int
	Last 		int
	Size		int
	Iterator	MessageQueueIterator
}


func NewIterator(messageQueue *MessageQueue) MessageQueueIterator{
	elementsLeft := messageQueue.Size
	if messageQueue.First <= messageQueue.Last{
		elementsLeft = (messageQueue.Last - messageQueue.First) + 1
	}
	return MessageQueueIterator{
		Queue: messageQueue,
		size: messageQueue.Size,
		currentIndex: messageQueue.First,
		elementsLeft: elementsLeft,
	}
}

func (i* MessageQueueIterator) HasNext() bool {
	fmt.Printf("current: %d, last: %d\n", i.currentIndex % i.size, i.Queue.Last)
	return i.elementsLeft > 0
}

func (i *MessageQueueIterator) Next() Message{
	fmt.Printf("grabbing element %d\n", i.currentIndex % i.size)
	current := i.Queue.Messages[i.currentIndex % i.size]
	i.currentIndex += 1
	i.elementsLeft -= 1
	return current
}


func New(size int) MessageQueue{
	return MessageQueue{
		Last: -1, 
		First: 0, 
		Size: size, 
		Messages: make([]Message, size), 
		Iterator: MessageQueueIterator{currentIndex: 0, size: size}}
}


func (mq *MessageQueue) AddMessage(message Message){
	// Keep track of the earliest message within the queue
	mq.Messages[(mq.Last + 1) % mq.Size] = message
	mq.increment()

}

func (mq *MessageQueue) increment(){
	if mq.Last == -1{
		mq.Last += 1
		return
	}

	if mq.Last == mq.Size - 1  && mq.First == 0{
		mq.incrementFirst()
		mq.incrementLast()
		return
	}

	mq.incrementLast()
	if mq.Last == mq.First {
		mq.incrementFirst()
	}
}


func (mq *MessageQueue) incrementFirst(){
	if mq.First == mq.Size - 1{
		mq.First = 0

	}else { mq.First += 1 }
}

func (mq *MessageQueue) incrementLast(){
	if mq.Last == mq.Size - 1 {
		mq.Last = 0
	} else { mq.Last += 1}
}







