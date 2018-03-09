package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"log"
)

var (
	rustedIronURI = "http://localhost:8000/queue/1"
)

type Msg struct {
	Body string `json:"body"`
}

type Reservation struct {
	N int `json:"n"`
}

type RustedIronRunner struct{}

func (rir *RustedIronRunner) setupQueues(queues []string) {}

func (rir *RustedIronRunner) Name() string { return "RustedIronMQ" }

func (rir *RustedIronRunner) Produce(name, body string, messages int) {
	produceURI := rustedIronURI
	produceURI += "/messages"

	msgs := make([]Msg, messages)
	for i := 0; i < messages; i++ {
		msgs[i] = Msg{body}
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msgs)

	_, err := http.Post(
		produceURI,
		"application/json; charset=utf-8",
		b)
	if err != nil {
		log.Println(err);
	}
}

func (rir * RustedIronRunner) Consume(name string, messages int) {
	reserveURI := rustedIronURI
	reserveURI += "/reservations"
	reservationBody := Reservation{messages}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reservationBody)

	_, err := http.Post(
		reserveURI,
		"application/json; charset=utf-8",
		b)
	if err != nil {
		log.Println(err);
	}
}
