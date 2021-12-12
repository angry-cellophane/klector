package api

import "net/http"

func Create() {
	http.HandleFunc("/api/v1/event")
}
