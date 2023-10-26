package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const CORELATION_ID_KEY = "corelation_id"

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func loggerMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Incoming request:", r.RequestURI)
		h(w, r)
	}
}

func requestTimingMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		h(w, r)
		after := time.Now()

		tookMs := after.Sub(before).Microseconds()
		log.Println("Request took in microseconds:", tookMs)
	}
}

func addRequestId(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		correlationId := uuid.New().String()

		ctx := context.WithValue(r.Context(), CORELATION_ID_KEY, correlationId)

		r = r.WithContext(ctx)

		h(w, r)
	}
}

func applyMiddlewares(h http.HandlerFunc, mm []MiddlewareFunc) http.HandlerFunc {
	for _, m := range mm {
		h = m(h)
	}

	return h
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world\n")

	val := r.Context().Value(CORELATION_ID_KEY)
	corelationId, _ := val.(string)

	log.Println("Request with id", corelationId)
}

func main() {
	middlewares := []MiddlewareFunc{addRequestId, loggerMiddleware, requestTimingMiddleware}

	handlerWithMiddlewares := applyMiddlewares(testHandler, middlewares)

	http.HandleFunc("/test", handlerWithMiddlewares)

	http.ListenAndServe(":8080", nil)
}
