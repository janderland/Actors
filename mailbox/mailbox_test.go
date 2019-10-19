package mailbox

import (
	"context"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"
)

const size = 10

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	os.Exit(m.Run())
}

func TestNewMailBox(t *testing.T) {
	box := NewMailBox(size)

	priorities := make([]int, size)
	for i := range priorities {
		priorities[i] = rand.Int()
	}

	t.Logf("Priorities: %v", priorities)

	for _, p := range priorities {
		box.Put(p, p)
	}

	sort.Sort(sort.IntSlice(priorities))

	for range priorities {
		result := box.Get(context.Background())
		t.Log(result)
		/*
		if *result.(*int) != p {
			t.Errorf("not equal")
		}
		 */
	}
}
