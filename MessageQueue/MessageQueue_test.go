package messagequeue

import (
	"fmt"
	"testing"
)


func TestMakeNewQueue(t *testing.T){
	messages := New(10)
	if messages.Size != 10 {
		t.Fatalf("size initiated incorrectly")
	}
}

func TestAddMessage(t *testing.T){
	messages := New(10)
	messages.AddMessage(Message{RawText: "Hello"})
	if messages.First != 0{
		t.Fatalf("Error first message should be at index 0")
	}
	if messages.Last != 0{
		t.Fatalf("Error last message should be at index 1")
	}
}

func TestMaxMessages(t *testing.T){
	messages := New(4)
	tests := []struct {
		E_first int
		E_last int
	} {
		{E_first: 0, E_last: 0},
		{E_first: 0, E_last: 1},
		{E_first: 0, E_last: 2},
		{E_first: 0, E_last: 3},
	}
	for _ , test := range tests{
		messages.AddMessage(Message{RawText: "Hello"})
		if messages.First != test.E_first{
			t.Fatalf("Error first expected %d got %d", test.E_first, messages.First )
		}
		if messages.Last != test.E_last{
			t.Fatalf("Error last expected %d got %d", test.E_last, messages.Last )
		}
	}
}


func TestDoubleRollOver( t *testing.T){
	messages := New(4)

	tests := []struct {
		E_first int
		E_last int
	} {
		{E_first: 0, E_last: 0},
		{E_first: 0, E_last: 1},
		{E_first: 0, E_last: 2},
		{E_first: 0, E_last: 3},
		{E_first: 1, E_last: 0},
		{E_first: 2, E_last: 1},
		{E_first: 3, E_last: 2},
		{E_first: 0, E_last: 3},

	}
	for _ , test := range tests{
		messages.AddMessage(Message{RawText: "Hello"})

		if messages.First != test.E_first{
			t.Fatalf("Error first expected %d got %d", test.E_first, messages.First )
		}
		if messages.Last != test.E_last{
			t.Fatalf("Error last expected %d got %d", test.E_last, messages.Last )
		}
	}
}

func TestMessageIterator( t *testing.T ) {
	messages :=  New(4)
	tests := []string{
		"1","2","3","4",
	}

	messages.AddMessage(Message{RawText: "1"})
	messages.AddMessage(Message{RawText: "2"})
	messages.AddMessage(Message{RawText: "3"})
	messages.AddMessage(Message{RawText: "4"})

	iterator := NewIterator(&messages)
	for _, test := range tests {
		message := iterator.Next()
		if message.RawText != test{
			t.Fatalf("error expected %s got %s", message.RawText, test)
		}
		fmt.Println(message.RawText)
	}

	messages.AddMessage(Message{RawText: "5"})
	messages.AddMessage(Message{RawText: "6"})

	tests = []string{
		"3","4","5","6",
	}

	iterator = NewIterator(&messages)
	for _, test := range tests {
		message := iterator.Next()
		if message.RawText != test{
			t.Fatalf("error expected %s got %s", test, message.RawText)
		}
		fmt.Println(message.RawText)
	}
	messages.AddMessage(Message{RawText: "7"})
	messages.AddMessage(Message{RawText: "8"})
	tests = []string{
		"5","6","7","8",
	}
	iterator = NewIterator(&messages)
	for _, test := range tests {
		message := iterator.Next()
		if message.RawText != test{
			t.Fatalf("error expected %s got %s", test, message.RawText)
		}
		fmt.Println(message.RawText)
	}
}

