package router

import (
	"context"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/lnquy/eventsourcing/config"
	"github.com/lnquy/eventsourcing/datastore"
	"github.com/lnquy/eventsourcing/model"
	"github.com/rs/xid"
	"net/http"
)

var (
	cfg *config.Config = config.LoadEnvConfig()
	ps  *datastore.PersonStorage
)

func init() {
	cfg = config.LoadEnvConfig()
	var err error
	ps, err = datastore.NewStorage(context.Background(), cfg.DataStore)
	if err != nil {
		logrus.Fatalf("router: Failed to create datastore client")
	}
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	p := &model.Person{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p.Id = xid.New().String() // Randomize Id

	//cEvent := &event.PersonCreated{
	//	Person: *p,
	//}
	// TODO: Save event
	savedP, err := ps.Save(context.Background(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render.JSON(w, r, savedP)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	getP, err := ps.Get(context.Background(), pId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, getP)
}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}

	p := &model.Person{}
	if err := json.NewDecoder(r.Body).Decode(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedP, err := ps.Save(context.Background(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render.JSON(w, r, updatedP)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	pId := chi.URLParam(r, "pid")
	if pId == "" {
		http.Error(w, "Invalid PersonID", http.StatusBadRequest)
		return
	}
	err := ps.Delete(context.Background(), pId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	render.JSON(w, r, "Ok")
}
