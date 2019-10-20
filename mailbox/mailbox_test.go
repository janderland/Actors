package mailbox

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)

var seed *int64

func TestMain(m *testing.M) {
	const seedFlagName = "seed"
	const seedFlagUsage = "seeds randomness"
	seedFlagDefault := time.Now().UTC().UnixNano()

	seed = flag.Int64(seedFlagName, seedFlagDefault, seedFlagUsage)
	rand.Seed(*seed)

	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	const mailBoxSize = 10
	box := NewMailBox(mailBoxSize)

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	message := box.Get(ctx)
	t.Logf("got: '%v'", message)

	if message != nil {
		t.Fatalf("'%v' != 'nil'", message)
	}
}

// Test basic operation in `Array` mode.
func TestPutThenGet(t *testing.T) {
	t.Logf("seed: %d", *seed)

	const mailBoxSize = 5
	const message = "msg"
	const priority = 0

	box := NewMailBox(mailBoxSize)

	box.Put(Message{message, priority})
	t.Logf("put: '%v'", message)

	received := box.Get(context.Background())
	t.Logf("got: '%v'", received)

	if received != message {
		t.Fatalf("'%v' != '%v'", message, received)
	}
}

// Test basic operation in `Array` mode.
func TestGetThenPut(t *testing.T) {
	t.Logf("seed: %d", *seed)

	const mailBoxSize = 5
	const message = "msg"
	const priority = 0

	box := NewMailBox(mailBoxSize)

	receivedCh := make(chan interface{})
	go func() {
		received := box.Get(context.Background())
		t.Logf("got: '%v'", received)
		receivedCh <- received
	}()

	box.Put(Message{message, priority})
	t.Logf("put: '%v'", message)

	if received := <-receivedCh; received != message {
		t.Fatalf("'%v' != '%v'", message, received)
	}
}

// Test basic operation in `Heap` mode.
func TestPutsThenGets(t *testing.T) {
	t.Logf("seed: %d", *seed)

	const numOfMessages = 5

	box := NewMailBox(numOfMessages)

	// Messages are defined as an int. Their contents are also used as their
	// priority.
	messages := make([]int, numOfMessages)
	for i := range messages {
		messages[i] = rand.Int()
	}

	for _, p := range messages {
		box.Put(Message{p, p})
		t.Logf("put: '%v'", p)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(messages)))

	for i := range messages {
		received := box.Get(context.Background())
		t.Logf("got: '%v'", received)
		if received != messages[i] {
			t.Fatalf("'%v' != '%v'", messages[i], received)
		}
	}
}

func TestConcurrently(t *testing.T) {
	t.Logf("seed: %d", *seed)

	const maxNumOfMessages = 50
	numOfMessages := rand.Intn(maxNumOfMessages) + 1
	t.Logf("Num Of Messages: %d", numOfMessages)

	const maxNumOfGetters = 10
	numOfGetter := rand.Intn(maxNumOfGetters) + 1
	t.Logf("Num Of Getters: %d", numOfGetter)

	const maxNumOfPutters = 10
	numOfPutters := rand.Intn(maxNumOfPutters) + 1
	t.Logf("Num Of Putters: %d", numOfPutters)

	box := NewMailBox(numOfMessages)

	// Messages are defined as an int. Their contents are also used as their
	// priority.
	messages := make([]int, numOfMessages)
	for i := range messages {
		messages[i] = rand.Int()
	}

	// Waits for all getters to return.
	wg := sync.WaitGroup{}
	wg.Add(numOfGetter)

	// Tells getters to return.
	ctx, stopGetting := context.WithCancel(context.Background())

	// Used by getters to forward messages to main go-routine.
	receivedMsgCh := make(chan int, numOfMessages)

	// Start go-routines that get messages.
	for i := 0; i < numOfGetter; i++ {
		go func() {
			defer wg.Done()
			for {
				message := box.Get(ctx)
				if message == nil {
					return
				}
				t.Logf("got: '%v'", message)
				receivedMsgCh <- message.(int)
			}
		}()
	}

	// Used by putters to request another message to put.
	requestMsgCh := make(chan chan int)

	// Start go-routine to provide the putters with messages.
	go func() {
		for i := range messages {
			respCh := <-requestMsgCh
			respCh <- messages[i]
		}
	}()

	// Setup go-routines that put messages.
	for i := 0; i < numOfPutters; i++ {
		go func() {
			for {
				ch := make(chan int)
				requestMsgCh <- ch
				message := <-ch
				box.Put(Message{message, message})
				t.Logf("put: '%v'", message)
			}
		}()
	}

	// Wait for each message to be received once.
	for range messages {
		received := <-receivedMsgCh

		found := false
		for i := range messages {
			if messages[i] == received {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("invalid receive: '%v'", received)
		}
	}

	// Shut down go-routines.
	stopGetting()
	wg.Wait()
}
