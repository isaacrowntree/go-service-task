package parser

import (
	"sort"
	"time"

	"github.com/isaacrowntree/go-service-task/pkg/reader"
	"github.com/isaacrowntree/go-service-task/pkg/structs"
)

const CsvColumns = 3

func ParseTimeStamp(date string) time.Time {
	layout := "2006-01-02T15:04:05Z"
	timeStamp, _ := time.Parse(layout, date)

	return timeStamp
}

func checkBounds(records [][]string, p structs.Parameters) bool {
	firstTime := ParseTimeStamp(records[0][0])
	lastTime := ParseTimeStamp(records[len(records)-1][0])

	return (firstTime.Before(p.To) && lastTime.After(p.From) ||
		firstTime.After(p.To) && firstTime.Before(p.From) ||
		lastTime.Before(p.From) && lastTime.After(p.To))
}

func getRecords(records [][]string, p structs.Parameters) []structs.LogRecord {
	var collected []structs.LogRecord
	for _, v := range records {
		if len(v) == CsvColumns {
			r := structs.LogRecord{
				EventTime: ParseTimeStamp(v[0]),
				Email:     v[1],
				SessionId: v[2],
			}
			if p.To.Before(r.EventTime) && p.From.After(r.EventTime) {
				collected = append(collected, r)
			}
		}
	}
	return collected
}

func checkFiles(files []string, p structs.Parameters, done chan []structs.LogRecord) {
	var found []structs.LogRecord
	for i := 0; i < len(files); i++ {
		records := reader.ReadCsvFile(files[i])
		dataExists := checkBounds(records, p)
		if dataExists {
			found = append(found, getRecords(records, p)...)
		}
	}
	done <- found
}

func sortData(data []structs.LogRecord) []structs.LogRecord {
	sort.Slice(data[:], func(i, j int) bool {
		return data[i].EventTime.Before(data[j].EventTime)
	})
	return data
}

func ParseFile(chunks [][]string, p structs.Parameters) []structs.LogRecord {
	var chunkSize = len(chunks)
	channels := make([]chan []structs.LogRecord, chunkSize)

	for i := 0; i < chunkSize; i++ {
		channels[i] = make(chan []structs.LogRecord)
		go checkFiles(chunks[i], p, channels[i])
	}

	returnMap := make(map[int][]structs.LogRecord)
	for i := 0; i < chunkSize; i++ {
		returnMap[i] = <-channels[i]
	}

	var output []structs.LogRecord
	for _, v := range returnMap {
		output = append(output, v...)
	}

	return sortData(output)
}
