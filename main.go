package main

import (
	"encoding/json"
	"fmt"
	"go-temperature-rest-example/temperaturestore"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	server := NewTemperatureServer()
	mux.HandleFunc("/temperature/", server.temperatureHandler)
	mux.HandleFunc("/city/", server.cityHandler)
	mux.HandleFunc("/datetime/", server.datetimeHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}

type TemperatureServer struct {
	store *temperaturestore.TemperatureStore
}

func NewTemperatureServer() *TemperatureServer {
	store := temperaturestore.New()
	return &TemperatureServer{store: store}
}

func (ts *TemperatureServer) temperatureHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/temperature/" {
		if req.Method == http.MethodPost {
			ts.createTemperatureHandler(w, req)
		} else if req.Method == http.MethodGet {
			ts.getAllTemperaturesHandler(w, req)
		} else if req.Method == http.MethodDelete {
			ts.deleteAllTemperaturesHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("use method GET, POST or DELETE at /temperature/, don`t %v", req.Method), http.StatusMethodNotAllowed)
		}
	} else {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")
		if len(pathParts) < 2 {
			http.Error(w, "expect /temperature/<id> in temperature handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			ts.deleteTemperatureHandler(w, req, int(id))
		} else if req.Method == http.MethodGet {
			ts.getTemperatureHandler(w, req, int(id))
		} else {
			http.Error(w, fmt.Sprintf("use method GET< POST or DELETE at /temperature/<id>, don`t %v", req.Method), http.StatusMethodNotAllowed)
		}
	}
}

func (ts *TemperatureServer) createTemperatureHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling temperature create at %s\n", req.URL.Path)

	type RequestTemperature struct {
		Value    int       `json:"value"`
		City     string    `json:"city"`
		Datetime time.Time `json:"datetime"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "use application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rt RequestTemperature
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTemperature(rt.Value, rt.City, rt.Datetime)
	js, err := json.Marshal(ResponseId{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (ts *TemperatureServer) getAllTemperaturesHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all temperature at %s\n", req.URL.Path)

	allTemperatures := ts.store.GetAllTemperatures()
	js, err := json.Marshal(allTemperatures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (ts *TemperatureServer) deleteAllTemperaturesHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all temperatures at %s\n", req.URL.Path)
	ts.store.DeleteAllTemperatures()
}

func (ts *TemperatureServer) deleteTemperatureHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handing delete task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTemperature(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (ts *TemperatureServer) getTemperatureHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get task at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	temperature, err := ts.store.GetTemperature(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(temperature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (ts *TemperatureServer) cityHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling temperatures by city at %s\n", req.URL.Path)

	city := req.PathValue("city")

	temperatures := ts.store.GetTemperaturesByCity(city)
	js, err := json.Marshal(temperatures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

func (ts *TemperatureServer) datetimeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling temperatures by datetime at %s\n", req.URL.Path)

	badRequestError := func() {
		http.Error(w, fmt.Sprintf("use /datetime/<year>/<month>/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		badRequestError()
		return
	}

	temperatures := ts.store.GetTemperaturesByDatetime(year, time.Month(month), day)
	js, err := json.Marshal(temperatures)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
