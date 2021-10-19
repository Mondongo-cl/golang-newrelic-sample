package main

import (
	"bltrain/heathcheck/testslogic/kafka/producer"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(producer.ProducerPostHandler)
	http.Handle("/kafka", handler)
	http.ListenAndServe(":8080", handler)
}
