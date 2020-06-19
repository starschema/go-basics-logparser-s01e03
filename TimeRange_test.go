package logparser

import (
	"testing"
	"time"
)

func TestTimeRange_Includes(t *testing.T) {

	data :=  []struct {
		startOffset, endOffset int
		otherTimeOffset int
		want bool
	}{
		{ 0, 5, 7, false},
		{ 0, 5, 3, true},
		{ 0, 5, -3, false},
	}

	for _, test := range data {
		tr := TimeRange{
			time.Now().Add( time.Duration(test.startOffset) * time.Minute),
			time.Now().Add( time.Duration(test.endOffset) * time.Minute) ,
		}
		otherTime := time.Now().Add(time.Duration(test.otherTimeOffset) * time.Minute)

		if tr.Includes(otherTime) != test.want {
			t.Fatalf("TimeRange %v  include test resulted in %v with argument %v", tr, test.want, otherTime)
		}
	}



}

func TestTimeRange_Intersects(t *testing.T) {

	data :=  []struct {
		startOffset, endOffset int
		otherStarOffset, otherEndOffset int
		want bool
	}{
		{ 0, 5, 6, 7, false},
		{ 0, 5, 3, 7, true},
		{ 0, 5, -3, -1, false},
	}

	for _, test := range data {
		tr := TimeRange{
			time.Now().Add( time.Duration(test.startOffset) * time.Minute),
			time.Now().Add( time.Duration(test.endOffset) * time.Minute) ,
		}
		other := TimeRange{
			time.Now().Add( time.Duration(test.otherStarOffset) * time.Minute),
			time.Now().Add( time.Duration(test.otherEndOffset) * time.Minute) ,
		}

		if tr.Intersects(other) != test.want {
			t.Fatalf("TimeRange %v  include test resulted in %v with argument %v", tr, test.want, other)
		}
	}



}
