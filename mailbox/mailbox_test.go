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

func TestNewMailBox(t *testing.T) {
	const numOfMessages = 5

	box := NewMailBox(numOfMessages)

	// Messages are defined as int. Their value is also used as the priority.
	messages := make([]int, numOfMessages)
	for i := range messages {
		messages[i] = rand.Int()
	}

	for _, p := range messages {
		box.Put(p, p)
	}

	sort.Sort(sort.IntSlice(messages))
	t.Logf("Priorities: %v", messages)

	for range messages {
		result := box.Get(context.Background())
		t.Log(result)
	}
}
