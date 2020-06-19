package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"logparser"
	"net/http"
	"time"
)


func createService() chan<- logparser.FindRequest {

	logFileGroups := []logparser.LogFileGroup {
		{ "/Users/gyulalaszlo/Documents/TableauLogs/logs/httpd/access.*.log", "access"},
	}

	parsers := []logparser.LogParser {
		&logparser.DummyLogParser{},
		&logparser.AccesLogParser{},
	}

	serviceChan := logparser.MakeLogProcessorService(parsers,logFileGroups)


	return serviceChan
}


type FindJsonRequest struct {
	Start string `json:"start"`
	End string `json:"end"`
}


func main() {

	serviceChan := createService()
	//start, _ := time.Parse("2006-01-02T15:04:05", "2020-06-11T13:48:00.000")
	//end, _ := time.Parse("2006-01-02T15:04:05", "2020-06-12T13:48:00.000")



	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {

		var req FindJsonRequest

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}


		start, err := time.Parse("2006-01-02T15:04:05", req.Start)
		if err != nil {
			http.Error(w, "Malformed start time", http.StatusBadRequest)
			return
		}

		end, err := time.Parse("2006-01-02T15:04:05", req.End)
		if err != nil {
			http.Error(w, "Malformed end time", http.StatusBadRequest)
			return
		}


		lines, err := logparser.FindWithService(serviceChan, logparser.TimeRange{ start, end })

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonLinesBytes, err := json.Marshal(lines)


		w.Write(jsonLinesBytes)



		//fmt.Fprintf(w, "Hello, %q %v", html.EscapeString(r.URL.Path), lines)
	})

	logrus.Info("Runnig server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}