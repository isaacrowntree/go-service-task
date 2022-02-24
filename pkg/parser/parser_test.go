package parser

import (
	"testing"

	"github.com/isaacrowntree/go-service-task/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestParseTimeStamp(t *testing.T) {
	t.Run("when the timestamp is expected", func(t *testing.T) {
		timeStamp := ParseTimeStamp("2001-02-20T23:14:31Z")

		assert.Equal(t, "2001-02-20 23:14:31 +0000 UTC", timeStamp.String())
	})

	t.Run("when the timestamp is a strange value it returns a default", func(t *testing.T) {
		timeStamp := ParseTimeStamp("is-not-a-date")

		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", timeStamp.String())
	})
}

func TestCheckBounds(t *testing.T) {

	var records = [][]string{
		{"2004-02-02T22:45:20Z", "user1@test.com", "fake-guid-1"},
		{"2004-02-02T22:46:20Z", "user2@test.com", "fake-guid-2"},
		{"2004-02-02T22:47:20Z", "user3@test.com", "fake-guid-3"},
		{"2005-02-02T22:47:20Z", "user4@test.com", "fake-guid-4"},
	}

	t.Run("when the bounds don't match at all", func(t *testing.T) {
		params := structs.Parameters{
			From: ParseTimeStamp("2021-07-06T23:00:00Z"),
			To:   ParseTimeStamp("2021-07-06T23:00:00Z"),
		}

		boundsResult := checkBounds(records, params)

		assert.Equal(t, false, boundsResult)
	})

	t.Run("when the 'to' value is within the range", func(t *testing.T) {
		params := structs.Parameters{
			From: ParseTimeStamp("2021-07-06T23:00:00Z"),
			To:   ParseTimeStamp("2004-03-06T23:00:00Z"),
		}

		boundsResult := checkBounds(records, params)

		assert.Equal(t, true, boundsResult)
	})

	t.Run("when the 'from' value is within the range", func(t *testing.T) {
		params := structs.Parameters{
			From: ParseTimeStamp("2005-01-03T23:00:00Z"),
			To:   ParseTimeStamp("2003-04-06T23:00:00Z"),
		}

		boundsResult := checkBounds(records, params)

		assert.Equal(t, true, boundsResult)
	})

	t.Run("when both values are in the range", func(t *testing.T) {
		params := structs.Parameters{
			From: ParseTimeStamp("2005-01-03T23:00:00Z"),
			To:   ParseTimeStamp("2004-02-03T23:00:00Z"),
		}

		boundsResult := checkBounds(records, params)

		assert.Equal(t, true, boundsResult)
	})

	t.Run("when both values include the range", func(t *testing.T) {
		params := structs.Parameters{
			From: ParseTimeStamp("2003-01-03T23:00:00Z"),
			To:   ParseTimeStamp("2006-02-03T23:00:00Z"),
		}

		boundsResult := checkBounds(records, params)

		assert.Equal(t, true, boundsResult)
	})
}

func TestSortData(t *testing.T) {
	var records = []structs.LogRecord{
		{
			EventTime: ParseTimeStamp("2005-02-02T22:47:20Z"),
			Email:     "user4@test.com",
			SessionId: "fake-guid-4",
		},
		{
			EventTime: ParseTimeStamp("2004-02-02T22:45:20Z"),
			Email:     "user1@test.com",
			SessionId: "fake-guid-1",
		},
		{
			EventTime: ParseTimeStamp("2004-02-02T22:47:20Z"),
			Email:     "user3@test.com",
			SessionId: "fake-guid-3",
		},
		{
			EventTime: ParseTimeStamp("2004-02-02T22:46:20Z"),
			Email:     "user2@test.com",
			SessionId: "fake-guid-2",
		},
	}

	output := sortData(records)

	assert.Equal(t, []structs.LogRecord{
		{
			EventTime: ParseTimeStamp("2004-02-02T22:45:20Z"),
			Email:     "user1@test.com",
			SessionId: "fake-guid-1",
		},
		{
			EventTime: ParseTimeStamp("2004-02-02T22:46:20Z"),
			Email:     "user2@test.com",
			SessionId: "fake-guid-2",
		},
		{
			EventTime: ParseTimeStamp("2004-02-02T22:47:20Z"),
			Email:     "user3@test.com",
			SessionId: "fake-guid-3",
		},
		{
			EventTime: ParseTimeStamp("2005-02-02T22:47:20Z"),
			Email:     "user4@test.com",
			SessionId: "fake-guid-4",
		},
	}, output)
}
