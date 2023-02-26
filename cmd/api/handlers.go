package main

import (
	"broker-service/event"
	"broker-service/notes"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReqPayload struct {
	Action string `json:"action,omitempty"`
	Name   string `json:"name,omitempty"`
	NoteID int    `json:"note_id,string,omitempty"`
	Data   Note   `json:"data,omitempty"`
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var reqPayload ReqPayload
	err := app.readJSON(w, r, &reqPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send payload to the notes service
	app.ConnectToNotesService(w, reqPayload)
}

// ConnectToNotesService base on action from the payload calls propper function
func (app *Config) ConnectToNotesService(w http.ResponseWriter, reqPayload ReqPayload) {
	switch reqPayload.Action {
	case "get-notes-list":
		// app.getNotesList(w)
		app.getNotesViaRPC(w, reqPayload)
	case "get-note-by-id":
		// app.getNoteByID(w, reqPayload)
		app.getNotesViaRPC(w, reqPayload)
	case "add-note", "update-note", "delete-note":
		app.sendEvent(w, reqPayload)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

// getNotesList sends request to the notes-service and returns list of notes in JSON format
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

	// make sure the status is accepted
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling notes service"))
		return
	}

	var jsonFromNotesService jsonResponse

	// read response from the notes-service
	err = json.NewDecoder(response.Body).Decode(&jsonFromNotesService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, jsonFromNotesService)
}

// getNoteByID sends request to the notes-service and returns note's data by given id in JSON format
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

// sendEvent handles methods that modify data.
// Sends events with Notes payload to RabbitMQ, then listner-sevice take the payload from RabiitMQ,
// and following payload's action sends request to the notes'service.
func (app *Config) sendEvent(w http.ResponseWriter, payload ReqPayload) {
	err := app.pushToQueue(payload, "note-info")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payloadResponse jsonResponse

	payloadResponse.Error = false
	payloadResponse.Message = "Operation's been successfully accomplished."

	app.writeJSON(w, http.StatusAccepted, payloadResponse)
}

// pushToQueue publishes events to RabbitMQ.
func (app *Config) pushToQueue(payload ReqPayload, msg string) error {

	// create producer
	producer, err := event.NewProducer(app.Rabbit)

	if err != nil {
		return err
	}

	jsonData, _ := json.Marshal(&payload)
	err = producer.Push(string(jsonData), "note-message")

	if err != nil {
		return err
	}

	return nil
}

type RPCPayload struct {
	Action string
	ID     int
}

// getNotesViaRPC creates client and invokes action from RPCServer.
// Handles methods that getting data.
func (app *Config) getNotesViaRPC(w http.ResponseWriter, requestPayload ReqPayload) {
	// connect to the RPCServer
	client, err := rpc.Dial("tcp", "notes-service:5001")
	if err != nil {
		return
	}

	rpcPayload := RPCPayload{
		Action: requestPayload.Action,
		ID:     requestPayload.Data.ID,
	}

	var responseString string
	var response jsonResponse
	err = client.Call("RPCServer.ReadAction", rpcPayload, &responseString)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = json.Unmarshal([]byte(responseString), &response)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, response)
}

func (app *Config) getNoteByIDViaGRPC(w http.ResponseWriter, r *http.Request) {
	var reqPayload ReqPayload

	err := app.readJSON(w, r, &reqPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// create a grpc client connection
	conn, err := grpc.Dial("notes-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	defer conn.Close()

	c := notes.NewNoteServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var grpcResp *notes.NoteResponse

	grpcResp, err = c.GetNotes(ctx, &notes.NoteRequest{
		Payload: &notes.Payload{
			Id: int32(reqPayload.Data.ID),
		},
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var note Note = Note{
		ID:              int(grpcResp.Note.Id),
		Name:            grpcResp.Note.Name,
		Description:     grpcResp.Note.Description,
		TextColor:       grpcResp.Note.TextColor,
		BackgroundColor: grpcResp.Note.BackgroundColor,
	}

	var jsonResponse jsonResponse
	jsonResponse.Data = note
	jsonResponse.Error = false

	app.writeJSON(w, http.StatusAccepted, jsonResponse)
}
