package logparser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

type AccesLogParser struct {

}

func (a *AccesLogParser) TypeName() string {
	return "access"
}

func (a *AccesLogParser) Parse(filename string) ([]LogLine, error) {
	// open the file
	file, err := os.Open(filename) // For read access.
	if err != nil {
		return nil, fmt.Errorf("while opening '%v' for reading: %v", filename, err)
	}
	// close the file
	defer file.Close()

	// create the empty container
	logLines := make([]LogLine, 0)

	// read line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		lineParsedTime, err := ParseHttpLine(lineText)
		if err != nil {
			fmt.Println("ERROR:", err)
			continue
		}

		logLines = append(logLines, LogLine{
			TimeStamp: lineParsedTime,
			Text:      lineText,
			Filename:  filename,
		})


	}

	return logLines, nil
	panic("implement me")
}

var httpLogLineRe = regexp.MustCompile(`^.* - (\d+.\d+.\d+T\d+:\d+:\d+\.\d+)`)

func ParseHttpLine(line string) (time.Time, error) {

	result := httpLogLineRe.FindSubmatch([]byte(line))

	// if we have errors
	if result == nil {
		return time.Time{}, fmt.Errorf("cannot find the timestamp in the line '%v'", line)
	}

	// parse the time
	parsedTime, err := time.Parse("2006-01-02T15:04:05", string(result[1]))
	if err != nil {
		return time.Time{}, fmt.Errorf("while parsing time: %v", err)
	}

	return parsedTime, nil
}

