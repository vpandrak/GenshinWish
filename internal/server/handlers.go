package server

import (
	"io"
	"net/http"
)

func getUnit(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "getUnit")
	if err != nil {
		return
	}
}

func wishUnit() {

}

func Pay() {

}

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "working")
	if err != nil {
		return
	}

}
