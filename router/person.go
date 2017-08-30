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

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	p := &model.Person{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Id = xid.New().String() // Randomize Id

	epc := &event.PersonCreated{
		Id: p.Id,
		Name: p.Name,
		Age: p.Age,
	}
	data, err := json.Marshal(epc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	e := &event.Event{
		Id: xid.New().String(),
		Data: data,
		Type: event.PersonCreatedEvent,
		Timestamp: time.Now().String(),
	}

	model.SavePersonEvents(context.Background(), e, p.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render.JSON(w, r, p)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	p := &model.Person{}
	err := model.GetPersonAggregate(context.Background(), p, pId)

	//p, err := ps.Get(context.Background(), pId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, p)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	agg := &model.Person{}
	err := model.GetPersonAggregate(context.Background(), agg, pId)
	if agg.Id != pId {
		http.Error(w, "PersonID not existed", http.StatusBadRequest)
		return
	}

	p := &model.Person{
		Id: agg.Id,
	}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	epu := &event.PersonUpdated {
		Name: p.Name,
		Age: p.Age,
	}
	data, err := json.Marshal(epu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	e := &event.Event{
		Id: xid.New().String(),
		Data: data,
		Type: event.PersonUpdatedEvent,
		Timestamp: time.Now().String(),
	}
	model.SavePersonEvents(context.Background(), e, pId)
	render.JSON(w, r, p)
}
