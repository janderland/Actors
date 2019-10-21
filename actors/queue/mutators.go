package queue

import (
	"github.com/janderland/Actors/messages"
)

func NewMutators() Mutators {
	return Mutators{
		Conn: func(
			conn messages.Conn,
			old State,
		) (
			outs []interface{},
			new State,
		) {
			accept := make(messages.Accept, len(conn))
			copy(accept, conn)

			outs = []interface{}{&accept}
			new = old
			return
		},
	}
}
