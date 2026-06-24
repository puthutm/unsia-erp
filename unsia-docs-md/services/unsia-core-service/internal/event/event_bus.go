package event

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event types for the system
const (
	EventUserCreated    = "user.created"
	EventUserUpdated  = "user.updated"
	EventUserDeleted   = "user.deleted"
	EventUserLogin    = "user.login"
	EventUserLogout   = "user.logout"
	EventRoleAssigned = "role.assigned"
	EventRoleRevoked = "role.revoked"
)

// Event represents an event in the system
type Event struct {
	ID        string      `json:"id"`
	Type     string      `json:"type"`
	Source   string      `json:"source"`
	Data     interface{} `json:"data"`
	Metadata map[string]string `json:"metadata,omitempty"`
	UserID   *string     `json:"user_id,omitempty"`
	ActorID  *string     `json:"actor_id,omitempty"`
	Time     time.Time   `json:"time"`
}

// EventHandler is the interface for event handlers
type EventHandler interface {
	Handle(event Event) error
}

// EventBus is the event bus implementation
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish publishes an event
func (eb *EventBus) Publish(eventType string, data interface{}, opts ...EventOption) error {
	event := Event{
		ID:   uuid.New().String(),
		Type: eventType,
		Time: time.Now(),
		Data: data,
	}

	// Apply options
	for _, opt := range opts {
		opt(&event)
	}

	// Get handlers
	eb.mu.RLock()
	handlers := eb.handlers[eventType]
	eb.mu.RUnlock()

	// Execute handlers
	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

// EventOption is a functional option for events
type EventOption func(*Event)

// WithSource sets the event source
func WithSource(source string) EventOption {
	return func(e *Event) {
		e.Source = source
	}
}

// WithUserID sets the user ID
func WithUserID(userID string) EventOption {
	return func(e *Event) {
		e.UserID = &userID
	}
}

// WithActorID sets the actor ID
func WithActorID(actorID string) EventOption {
	return func(e *Event) {
		e.ActorID = &actorID
	}
}

// WithMetadata sets metadata
func WithMetadata(meta map[string]string) EventOption {
	return func(e *Event) {
		e.Metadata = meta
	}
}

// EventToJSON converts event to JSON
func EventToJSON(event Event) (string, error) {
	data, err := json.Marshal(event)
	return string(data), err
}

// JSONToEvent converts JSON to event
func JSONToEvent(jsonStr string) (Event, error) {
	var event Event
	err := json.Unmarshal([]byte(jsonStr), &event)
	return event, err
}

// EventStore stores events for replay
type EventStore struct {
	events []Event
	mu     sync.Mutex
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make([]Event, 0),
	}
}

// Append adds an event to the store
func (es *EventStore) Append(event Event) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.events = append(es.events, event)
	return nil
}

// GetAll returns all events
func (es *EventStore) GetAll() []Event {
	es.mu.Lock()
	defer es.mu.Unlock()
	events := make([]Event, len(es.events))
	copy(events, es.events)
	return events
}

// GetByType returns events by type
func (es *EventStore) GetByType(eventType string) []Event {
	es.mu.Lock()
	defer es.mu.Unlock()
	var result []Event
	for _, e := range es.events {
		if e.Type == eventType {
			result = append(result, e)
		}
	}
	return result
}

// Clear clears all events
func (es *EventStore) Clear() {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.events = make([]Event, 0)
}

// PrintfEventHandler logs events
type PrintfEventHandler struct {
	Log func(format string, args ...interface{})
}

func (h PrintfEventHandler) Handle(event Event) error {
	data, _ := json.Marshal(event.Data)
	h.Log("[EVENT] %s: %s", event.Type, string(data))
	return nil
}

// AsyncEventHandler handles events asynchronously
type AsyncEventHandler struct {
	handler EventHandler
	channel chan Event
}

func NewAsyncEventHandler(handler EventHandler, bufferSize int) *AsyncEventHandler {
	ch := make(chan Event, bufferSize)
	h := &AsyncEventHandler{
		handler: handler,
		channel: ch,
	}
	go h.process()
	return h
}

func (h *AsyncEventHandler) Handle(event Event) error {
	h.channel <- event
	return nil
}

func (h *AsyncEventHandler) process() {
	for event := range h.channel {
		if err := h.handler.Handle(event); err != nil {
			fmt.Printf("Error processing event: %v\n", err)
		}
	}
}

// GlobalEventBus is the global event bus
var globalEventBus = NewEventBus()

// Subscribe registers a global handler
func Subscribe(eventType string, handler EventHandler) {
	globalEventBus.Subscribe(eventType, handler)
}

// Publish publishes a global event
func Publish(eventType string, data interface{}, opts ...EventOption) error {
	return globalEventBus.Publish(eventType, data, opts...)
}
