package router

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/lnquy/eventsourcing/model"
	"github.com/rs/xid"
	"net/http"
	"github.com/lnquy/eventsourcing/event"
	"time"
)

// CreatePerson handler.
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	p := &model.Person{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.Id = xid.New().String() // Randomize person id

	// Create PersonCreated event
	epc := &event.PersonCreated{
		Id:   p.Id,
		Name: p.Name,
		Age:  p.Age,
	}
	// Encode/Marshal PersonCreated event to raw []byte data to store in the event.Event wrapper
	data, err := json.Marshal(epc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the event wrapper for PersonCreated event
	e := &event.Event{
		Id:        xid.New().String(),
		Data:      data,
		Type:      event.PersonCreatedEvent,
		Timestamp: time.Now().String(),
	}

	// Stores event to database
	model.SavePersonEvents(context.Background(), e, p.Id)

	// Response
	render.JSON(w, r, p)
}

// GetPerson handler.
func GetPerson(w http.ResponseWriter, r *http.Request) {
	// Verify personId parameter
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	// Build person aggregate form event logs
	p := &model.Person{}
	err := model.GetPersonAggregate(context.Background(), p, pId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	render.JSON(w, r, p)
}

// UpdatePerson handler.
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	// Verify personId parameter
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	// Verify if person existed or not
	agg := &model.Person{}
	err := model.GetPersonAggregate(context.Background(), agg, pId)
	if agg.Id != pId {
		http.Error(w, "PersonID not existed", http.StatusBadRequest)
		return
	}

	// Decode request body
	p := &model.Person{
		Id: agg.Id,
	}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create PersonUpdated event
	epu := &event.PersonUpdated{
		Name: p.Name,
		Age:  p.Age,
	}
	// Encode/Marshal PersonUpdated event to raw []byte data to store in the event.Event wrapper
	data, err := json.Marshal(epu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the event wrapper for PersonUpdated event
	e := &event.Event{
		Id:        xid.New().String(),
		Data:      data,
		Type:      event.PersonUpdatedEvent,
		Timestamp: time.Now().String(),
	}

	// Stores event to database
	model.SavePersonEvents(context.Background(), e, pId)

	// Response
	render.JSON(w, r, p)
}
