package main

import (
	"broker-service/event"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ReqPayload struct {
	Action string `json:"action,omitempty"`
	Name   string `json:"name,omitempty"`
	NoteID string `json:"note_id,omitempty"`
	Data   Note   `json:"data,omitempty"`
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload ReqPayload
	err := app.readJSON(w, r, &reqPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.ConnectToNotesService(w, reqPayload)
}

func (app *Config) ConnectToNotesService(w http.ResponseWriter, reqPayload ReqPayload) {
	switch reqPayload.Action {
	case "get-notes-list":
		app.getNotesList(w)
	case "get-note-by-id":
		app.getNoteByID(w, reqPayload)
		reqPayload.Name = "add"
		var note Note = Note{
			Name:            "test",
			Description:     "test descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest descriptiontest description",
			TextColor:       "#FDF6F2",
			BackgroundColor: "#EE8044",
		}

		reqPayload.Data = note
		app.sendEvent(w, reqPayload)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) getNotesList(w http.ResponseWriter) {
	request, err := http.NewRequest("GET", "http://notes-service", nil)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create http client
	client := &http.Client{}

	// send request
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling notes service"))
		return
	}

	var jsonFromNotesService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromNotesService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromNotesService)
}

func (app *Config) getNoteByID(w http.ResponseWriter, reqPayload ReqPayload) {
	jsonData, _ := json.Marshal(reqPayload)
	request, err := http.NewRequest("POST", "http://notes-service/get-note-by-id", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling notes service"))
		return
	}

	var jsonFromNotesService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromNotesService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromNotesService)
}
func (app *Config) sendEvent(w http.ResponseWriter, payload ReqPayload) {
	err := app.pushToQueue(payload, "test")
	if err != nil {
		log.Println(125, err)
		app.errorJSON(w, err)
		return
	}
	log.Println("poszło")
	var response jsonResponse
	response.Error = false
	response.Message = "ok"
	// app.writeJSON(w, http.StatusAccepted, response)
}

func (app *Config) pushToQueue(payload ReqPayload, msg string) error {

	producer, err := event.NewProducer(app.Rabbit)

	if err != nil {
		log.Println(140, err)
		return err
	}

	// payload := ReqPayload{
	// 	Name: name,
	// 	Data: msg,
	// }

	jsonData, _ := json.Marshal(&payload)
	err = producer.Push(string(jsonData), "add")

	if err != nil {
		log.Println(153, err)
		return err
	}

	return nil
}
