package mailbox

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	const seedFlagName = "seed"
	const seedFlagUsage = "seeds randomness"
	seedFlagDefault := time.Now().UTC().UnixNano()

	seed := flag.Int64(seedFlagName, seedFlagDefault, seedFlagUsage)
	rand.Seed(*seed)

	os.Exit(m.Run())
}

func TestPutAndGet(t *testing.T) {
	const mailBoxSize = 5
	const message = "m"
	const priority = 0

	box := NewMailBox(mailBoxSize)

	box.Put(message, priority)
	t.Logf("put: '%v'", message)

	received := box.Get(context.Background())
	t.Logf("got: '%v'", received)

	if received != message {
		t.Fatalf("'%v' != '%v'", message, received)
	}
}

func TestGetAndPut(t *testing.T) {
	const mailBoxSize = 5
	const message = "m"
	const priority = 0

	box := NewMailBox(mailBoxSize)

	receivedCh := make(chan interface{})
	go func() {
		received := box.Get(context.Background())
		t.Logf("got: '%v'", received)
		receivedCh <- received
	}()

	box.Put(message, priority)
	t.Logf("put: '%v'", message)

	if received := <-receivedCh; received != message {
		t.Fatalf("'%v' != '%v'", message, received)
	}
}

func TestPutsThenGets(t *testing.T) {
	const numOfMessages = 5

	box := NewMailBox(numOfMessages)

	// Messages are defined as int. Their value is also used as the priority.
	messages := make([]int, numOfMessages)
	for i := range messages {
		messages[i] = rand.Int()
	}

	for _, p := range messages {
		box.Put(p, p)
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
