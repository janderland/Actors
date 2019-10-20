package messages

import "github.com/janderland/Actors/mailbox"

// @Conn
type Conn []byte

const ConnPriority = 0

func (c Conn) Pack() mailbox.Message {
	return mailbox.Message{
		Contents: c,
		Priority: ConnPriority,
	}
}

// @Accept
type Accept []byte

const AcceptPriority = 0

func (a Accept) Pack() mailbox.Message {
	return mailbox.Message{
		Contents: a,
		Priority: AcceptPriority,
	}
}

// (@Conn @Chunk)Data
type Data struct {
	Conn  []byte
	Chunk []byte
}

const DataPriority = 0

func (d *Data) Pack() mailbox.Message {
	return mailbox.Message{
		Contents: d,
		Priority: DataPriority,
	}
}
