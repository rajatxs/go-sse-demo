package main

import "net/http"

func getServerMux() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/ping", sendPingResponse())
	mux.HandleFunc("/peers", sendPeerList())
	mux.HandleFunc("/send", sendMessage())
	mux.HandleFunc("/stream", startStreaming())
	return mux
}
