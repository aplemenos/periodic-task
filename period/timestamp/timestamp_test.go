package timestamp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamp_OneHour(t *testing.T) {
	t1, _ := time.Parse("20060102T150405Z", "20210714T204603Z")
	t2, _ := time.Parse("20060102T150405Z", "20210715T123456Z")

	oht := NewTimestamp(ONEHOUR)
	result := oht.GetMatchingTimestamps(t1, t2)

	expected := []string{
		"20210714T210000Z",
		"20210714T220000Z",
		"20210714T230000Z",
		"20210715T000000Z",
		"20210715T010000Z",
		"20210715T020000Z",
		"20210715T030000Z",
		"20210715T040000Z",
		"20210715T050000Z",
		"20210715T060000Z",
		"20210715T070000Z",
		"20210715T080000Z",
		"20210715T090000Z",
		"20210715T100000Z",
		"20210715T110000Z",
		"20210715T120000Z",
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d timestamps, but got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %s, but got %s", expected[i], result[i])
		}
	}
}

func TestTimestamp_OneDay(t *testing.T) {
	t1, _ := time.Parse("20060102T150405Z", "20211010T204603Z")
	t2, _ := time.Parse("20060102T150405Z", "20211115T123456Z")

	odt := NewTimestamp(ONEDAY)
	result := odt.GetMatchingTimestamps(t1, t2)

	expected := []string{
		"20211010T210000Z", "20211011T210000Z", "20211012T210000Z", "20211013T210000Z",
		"20211014T210000Z", "20211015T210000Z", "20211016T210000Z", "20211017T210000Z",
		"20211018T210000Z", "20211019T210000Z", "20211020T210000Z", "20211021T210000Z",
		"20211022T210000Z", "20211023T210000Z", "20211024T210000Z", "20211025T210000Z",
		"20211026T210000Z", "20211027T210000Z", "20211028T210000Z", "20211029T210000Z",
		"20211030T210000Z", "20211031T210000Z", "20211101T210000Z", "20211102T210000Z",
		"20211103T210000Z", "20211104T210000Z", "20211105T210000Z", "20211106T210000Z",
		"20211107T210000Z", "20211108T210000Z", "20211109T210000Z", "20211110T210000Z",
		"20211111T210000Z", "20211112T210000Z", "20211113T210000Z", "20211114T210000Z",
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d timestamps, but got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %s, but got %s", expected[i], result[i])
		}
	}
}

func TestTimestamp_OneMonth(t *testing.T) {
	t1, _ := time.Parse("20060102T150405Z", "20210214T204603Z")
	t2, _ := time.Parse("20060102T150405Z", "20211215T123456Z")

	omt := NewTimestamp(ONEMONTH)
	result := omt.GetMatchingTimestamps(t1, t2)

	expected := []string{
		"20210228T210000Z",
		"20210331T210000Z",
		"20210430T210000Z",
		"20210531T210000Z",
		"20210630T210000Z",
		"20210731T210000Z",
		"20210831T210000Z",
		"20210930T210000Z",
		"20211031T210000Z",
		"20211130T210000Z",
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d timestamps, but got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %s, but got %s", expected[i], result[i])
		}
	}
}

func TestTimestamp_OneYear(t *testing.T) {
	t1, _ := time.Parse("20060102T150405Z", "20180214T204603Z")
	t2, _ := time.Parse("20060102T150405Z", "20221115T123456Z")

	oyt := NewTimestamp(ONEYEAR)
	result := oyt.GetMatchingTimestamps(t1, t2)

	expected := []string{
		"20181231T210000Z",
		"20191231T210000Z",
		"20201231T210000Z",
		"20211231T210000Z",
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %d timestamps, but got %d", len(expected), len(result))
		return
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Expected %s, but got %s", expected[i], result[i])
		}
	}
}

func TestTimestamp_UnsupportedPeriod(t *testing.T) {
	owt := NewTimestamp("1w")
	assert.Nil(t, owt, "Unsupported period")
}
