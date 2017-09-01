/*
Entry point and server initialization for Twitch API
By: Connor Douthat
9/1/2017
*/
package main

import (
	"fmt"
	"strings"
	"errors"
	"encoding/json"
	"net/http"
	"./twitch"
)

func main() {
	fmt.Println("Booting the server...")

	// Serve the frontend index HTML at the root (and also 404s)
	http.HandleFunc("/", serveIndex)

	// API endpoints
	http.HandleFunc("/api/user", serveUser)
	http.HandleFunc("/api/channel", serveUserChannel)
	http.HandleFunc("/api/stream", serveUserStream)

	// Run your server
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(user)
}

func serveUserChannel(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(w, r)
	if err != nil {
		return
	}

	channel, err := twitch.GetChannelByUser(user.Id)
	if err != nil {
		// TODO: ideally, this should return more specific errors, and log unexpected ones
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(channel)
}

func serveUserStream(w http.ResponseWriter, r *http.Request) {
	user, err := getRequestUser(w, r)
	if err != nil {
		return
	}

	stream, err := twitch.GetStreamByUser(user.Id)
	if err != nil {
		// TODO: ideally, this should return more specific errors, and log unexpected ones
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(stream)
}

func getRequestUser(w http.ResponseWriter, r *http.Request) (user *twitch.TwitchUser, err error) {
	name, err := getRequestName(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	user, err = twitch.GetUserByName(name)
	if err != nil {
		// TODO: ideally, this should return more specific errors, and log unexpected ones
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func getRequestName(r *http.Request) (name string, err error) {
	r.ParseForm()

	formName := r.Form["name"]
	if len(formName) != 1 {
		err = errors.New("Requires name parameter")
		return
	}

	name = strings.TrimSpace(formName[0])
	if name == "" {
		err = errors.New("Name must not be empty")
		return
	}

	return
}
