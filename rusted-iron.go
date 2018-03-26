package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var (
	rustedIronURI = "http://localhost:8080"
)

type Msg struct {
	Body string `json:"body"`
}

type Reservation struct {
	N int `json:"n"`
	Delete bool `json:"delete"`
}

type RustedIronRunner struct{}

func (ir *RustedIronRunner) setupQueues(queues []string) {
	// qs, err := mq.List()
	// for _, q := range qs {
	// 	log.Println("[INFO] deleting queues")
	// 	err = q.Delete()
	// 	if err != nil {
	// 		log.Println("delete err", err)
	// 	}
	// }
	putQueueURI := rustedIronURI
	putQueueURI += "/queues"
	b := new(bytes.Buffer)

	for _, q := range queues {
		putURI := putQueueURI
		putURI += "/" + q
		req, err := http.NewRequest(http.MethodPut, putURI, b)
		req.Header.Add("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Println("err", err)
		}
	}
}

func (rir *RustedIronRunner) Name() string { return "RustedIronMQ" }

func (rir *RustedIronRunner) Produce(name, body string, messages int) {
	produceURI := rustedIronURI
	produceURI += "/queues/" + name
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
		log.Println(err)
	}
}

func (rir *RustedIronRunner) Consume(name string, messages int) {
	reserveURI := rustedIronURI
	reserveURI += "/queues/" + name
	reserveURI += "/reservations"
	reservationBody := Reservation{N: messages, Delete: true}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(reservationBody)

	_, err := http.Post(
		reserveURI,
		"application/json; charset=utf-8",
		b)
	if err != nil {
		log.Println(err)
	}
}
