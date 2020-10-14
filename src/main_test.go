package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestHandler(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			fmt.Fprintf(w, "127.0.0.1")
		}))
		defer ts.Close()

		qs := map[string]string{"name": "Joy"}
		res, err := handler(events.APIGatewayProxyRequest{QueryStringParameters: qs})
		if err != nil {
			t.Fatal("Everything should be ok")
		}

		expected := "{\"message\":\"Hello, Joy\"}"
		if res.Body != expected {
			t.Fatalf("Expected: %s, Actual: %s", expected, res.Body)
		}
	})
}
