package logparser

import "time"

type DummyLogParser struct {

}

func (d *DummyLogParser) TypeName() string {
	return "dummy"
}

func (d *DummyLogParser) Parse(filename string) ([]LogLine, error) {
	return []LogLine{
		{
			TimeStamp: time.Now(),
			Text:      "FOOBAR",
			Filename:  filename,
		},
	}, nil
}
