package logic

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jqiris/nano"
	"github.com/jqiris/nano/component"
	"github.com/jqiris/nano/examples/demo/tadpole/logic/protocol"
	"github.com/jqiris/nano/session"
)

// World contains all tadpoles
type World struct {
	component.Base
	*nano.Group
}

// NewWorld returns a world instance
func NewWorld() *World {
	return &World{
		Group: nano.NewGroup(uuid.New().String()),
	}
}

// Init initialize world component
func (w *World) Init() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		w.Leave(s)
		w.Broadcast("leave", &protocol.LeaveWorldResponse{ID: s.ID()})
		log.Println(fmt.Sprintf("session count: %d", w.Count()))
	})
}

// Enter was called when new guest enter
func (w *World) Enter(s *session.Session, msg []byte) error {
	w.Add(s)
	log.Println(fmt.Sprintf("session count: %d", w.Count()))
	return s.Response(&protocol.EnterWorldResponse{ID: s.ID()})
}

// Update refresh tadpole's position
func (w *World) Update(s *session.Session, msg []byte) error {
	return w.Broadcast("update", msg)
}

// Message handler was used to communicate with each other
func (w *World) Message(s *session.Session, msg *protocol.WorldMessage) error {
	msg.ID = s.ID()
	return w.Broadcast("message", msg)
}