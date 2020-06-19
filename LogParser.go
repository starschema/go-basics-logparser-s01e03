package logparser

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type LogLine struct {
	TimeStamp time.Time
	Text string
	Filename string
}

func (l LogLine)String() string {
	return fmt.Sprintf("(%v) : [%v] %v", l.Filename, l.TimeStamp, l.Text)
}


type LogParser interface {
	TypeName() string
	Parse(filename string) ([]LogLine, error)
}


type LogFileGroup struct {
	FileGlob string
	ParserName string
}

type indexedFile struct {
	filename string
	parser string

	timeRange TimeRange
}

type LogProcessor struct {
	parsers map[string]LogParser

	indexedFiles []indexedFile
}

func (l *LogProcessor)Index(groups []LogFileGroup) error {

	for _, group := range groups {
		matches, err := filepath.Glob(group.FileGlob)
		if err != nil {
			return fmt.Errorf("whilie attempting to glob for '%v': %v", group.FileGlob, err)
		}

		parser, ok := l.parsers[group.ParserName]
		if !ok {
			return fmt.Errorf("cannot find parser by name '%v'", group.ParserName)
		}

		for _, filename := range matches {
			logrus.WithFields(logrus.Fields{
				"filename": filename,
				"parser": parser.TypeName(),
			}).Infof("Processing file")

			logLines, err := parser.Parse(filename)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"filename": filename,
					"parser": parser.TypeName(),
					"error": err,
				}).Error("Error while parsing file")
				continue
			}

			if len(logLines) == 0 {
				logrus.WithFields(logrus.Fields{
					"filename": filename,
					"parser": parser.TypeName(),
				}).Warn("Zero log lines returned")
				continue
			}

			startTime := logLines[0].TimeStamp
			endTime := logLines[len(logLines) - 1].TimeStamp

			parsedFile := indexedFile{
				filename:  filename,
				parser:    parser.TypeName(),
				timeRange: TimeRange{startTime, endTime },
			}

			l.indexedFiles = append(l.indexedFiles, parsedFile)

			logrus.WithFields(logrus.Fields{
				"filename": filename,
				"fileCount": len(l.indexedFiles),
			}).Infof("Added file")
			//
		}
	}

	return nil
}

func (l *LogProcessor)FindLines(t TimeRange) ([]LogLine, error) {

	allLines := make([]LogLine, 0)

	for _, f := range l.indexedFiles {
		logrus.WithFields(logrus.Fields{
			"filename": f.filename,
			"timeRange": f.timeRange,
			"parser": f.parser,
		}).Info("Looking at file")

		if f.timeRange.Intersects(t) {

			logrus.WithFields(logrus.Fields{
				"filename": f.filename,
				"timeRange": f.timeRange,
				"parser": f.parser,
			}).Info("File matches time range")


			parser, ok := l.parsers[f.parser]
			if !ok {
				return nil, fmt.Errorf("cannot find parser by name '%v'", f.parser)
			}

			lines, err := parser.Parse(f.filename)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"filename": f.filename,
					"parser": parser.TypeName(),
					"error": err,
				}).Error("Error while parsing file")
				continue
			}

			for _, line := range lines {
				if t.Includes(line.TimeStamp) {
					allLines = append(allLines, line)
				}
			}
		}
	}

	return allLines, nil
}


func MakeLogProcessor(parsers []LogParser) *LogProcessor {
	parserMap := make(map[string]LogParser)
	for _, parser := range parsers {
		parserMap[parser.TypeName()] = parser
	}
	return &LogProcessor{
		parsers: parserMap,
		indexedFiles: make([]indexedFile, 0),
	}
}




type FindRequest struct {
	TimeRange TimeRange

	ResponseChan chan<- FindResponse
}

type FindResponse struct {
	Lines []LogLine
	Error error
}


func MakeLogProcessorService(parsers []LogParser, groups []LogFileGroup) chan<- FindRequest {
	requestsChan := make(chan FindRequest)

	go func(){
		lp := MakeLogProcessor(parsers)
		err := lp.Index(groups)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"fileGroups": groups,
				"error": err,
			}).Error("Error while indexing log file groups")

			close(requestsChan)
			return
		}

		for req := range requestsChan {
			lines, err := lp.FindLines(req.TimeRange)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"timerange": req.TimeRange,
				}).Error("Error while finding in time range")
			}

			response := FindResponse{
				Lines: lines,
				Error: err,
			}

			req.ResponseChan <- response
		}
	}()

	return requestsChan
}


func FindWithService(serviceChan chan<- FindRequest, timeRange TimeRange) ([]LogLine, error) {
	responseChan := make(chan FindResponse)
	serviceChan <- FindRequest{
		TimeRange:    timeRange,
		ResponseChan: responseChan,
	}

	response := <- responseChan

	return response.Lines, response.Error
}




