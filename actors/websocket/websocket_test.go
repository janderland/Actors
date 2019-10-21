package websocket

import (
	"bytes"
	"context"
	"flag"
	"github.com/janderland/Actors/mailbox"
	"github.com/janderland/Actors/messages"
	"math/rand"
	"os"
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

func TestAccept(t *testing.T) {
	t.Logf("seed: %d", *seed)

	const mailBoxSize = 5
	const maxStringSize = 10

	connBox := mailbox.NewMailBox(mailBoxSize)
	dataBox := mailbox.NewMailBox(mailBoxSize)

	size := rand.Intn(maxStringSize) + 1
	accept := messages.Accept(make([]byte, size))
	rand.Read(accept)

	size = rand.Intn(maxStringSize) + 1
	conn := messages.Conn(make([]byte, size))
	rand.Read(conn)

	data := messages.Data{}
	size = rand.Intn(maxStringSize) + 1
	data.Conn = make([]byte, size)
	rand.Read(data.Conn)
	size = rand.Intn(maxStringSize) + 1
	data.Chunk = make([]byte, size)
	rand.Read(data.Chunk)

	actor := Actor{
		State:   State{},
		MailBox: mailbox.NewMailBox(mailBoxSize),
		Peers: Peers{
			Conn: [1]mailbox.PeerMailBox{connBox},
			Data: [1]mailbox.PeerMailBox{dataBox},
		},
		Mutators: Mutators{
			Accept: func(in messages.Accept, old State) ([]interface{}, State) {
				if !bytes.Equal(in, accept) {
					t.Fatalf("'%v' != '%v'", in, accept)
				}

				return []interface{}{conn, data}, old
			},
		},
	}

	ctx, stop := context.WithCancel(context.Background())
	group := sync.WaitGroup{}
	group.Add(1)

	go actor.Act(ctx, &group)

	actor.MailBox.Put(accept.Pack())

	connGet := connBox.Get(context.Background())
	if !bytes.Equal(connGet.(messages.Conn), conn) {
		t.Fatalf("'%v' != '%v'", connGet, conn)
	}

	dataGet := dataBox.Get(context.Background())
	dataConn := dataGet.(*messages.Data).Conn
	if !bytes.Equal(dataConn, data.Conn) {
		t.Fatalf("'%v' != '%v'", dataConn, data.Conn)
	}
	dataChunk := dataGet.(*messages.Data).Chunk
	if !bytes.Equal(dataChunk, data.Chunk) {
		t.Fatalf("'%v' != '%v'", dataChunk, data.Chunk)
	}

	stop()
	group.Wait()
}
