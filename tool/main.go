package main

import (
	"fmt"
	"logparser"
	"time"
)

func main() {

	logFileGroups := []logparser.LogFileGroup {
		{ "/Users/gyulalaszlo/Documents/TableauLogs/logs/httpd/access.*.log", "access"},
	}

	parsers := []logparser.LogParser {
		&logparser.DummyLogParser{},
		&logparser.AccesLogParser{},
	}

	serviceChan := logparser.MakeLogProcessorService(parsers,logFileGroups)



	//lp := logparser.MakeLogProcessor(parsers)
	//
	//lp.Index(logFileGroups)
	//

	start, _ := time.Parse("2006-01-02T15:04:05", "2020-06-11T13:48:00.000")
	end, _ := time.Parse("2006-01-02T15:04:05", "2020-06-12T13:48:00.000")



	lines, err := logparser.FindWithService(serviceChan, logparser.TimeRange{start, end})
	if err != nil {
		panic(err)
	}


	fmt.Println(lines)
}
