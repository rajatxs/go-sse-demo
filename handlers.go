package main

import (
	"fmt"
	"log"
	"net/http"
)

// sends simple reply message
func sendPingResponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)
		sendResponse(w, 200, &ApiResponse{
			Message: "Pong!",
		})
	}
}

// sends list of active peers
func sendPeerList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []string = []string{}

		setCORS(w)

		if len(activePeers) > 0 {
			for p := range activePeers {
				list = append(list, p)
			}
		}

		sendResponse(w, 200, &ApiResponse{
			Message: "Peer list",
			Result:  list,
		})
	}
}

// accept join request from peer and start streaming
func startStreaming() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			id      string
			flusher http.Flusher
		)

		if r.URL.Query().Has("id") {
			id = r.URL.Query().Get("id")
		}

		if len(id) == 0 {
			sendResponse(w, 400, &ApiResponse{
				Message: "missing id query param",
			})
			return
		}

		setCORS(w)
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		if activePeers[id] == nil {
			activePeers[id] = &Peer{
				Id: id,
				Ch: make(chan string),
			}
			log.Printf("peer connected: %s\n", id)
		} else {
			sendResponse(w, 400, &ApiResponse{
				Message: "connection limit exceeds",
			})
			return
		}

		flusher, _ = w.(http.Flusher)
		flusher.Flush()

		for {
			var msg string

			flusher.Flush()

			select {
			case msg = <-activePeers[id].Ch:
				if _, err := fmt.Fprintf(w, "data: %s\n\n", msg); err != nil {
					log.Printf("couldn't write message %v\n", err)
				}
				log.Printf("message sent to %s\n", id)
				flusher.Flush()

			case <-r.Context().Done():
				defer r.Body.Close()
				close(activePeers[id].Ch)
				delete(activePeers, id)
				flusher.Flush()
				log.Printf("peer disconnected: %s\n", id)
				return
			}
		}
	}
}

// sending message to given active peer
func sendMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			rec string
			msg string
		)

		setCORS(w)

		if r.URL.Query().Has("to") {
			rec = r.URL.Query().Get("to")
		} else {
			sendResponse(w, 400, &ApiResponse{
				Message: "missing receiver id",
			})
			return
		}

		if r.URL.Query().Has("msg") {
			msg = r.URL.Query().Get("msg")
		} else {
			msg = ""
		}

		if len(rec) > 0 && activePeers[rec] != nil {
			log.Printf("sending message rec: %s, message: %s\n", rec, msg)
			activePeers[rec].Ch <- msg
		} else {
			sendResponse(w, 404, &ApiResponse{
				Message: "peer not connected",
			})
		}
	}
}
