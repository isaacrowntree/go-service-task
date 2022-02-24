package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/isaacrowntree/go-service-task/pkg/handler"
	"github.com/isaacrowntree/go-service-task/pkg/reader"
)

func TestMain(t *testing.T) {

	path := reader.GetCurrDir()
	sourceFile := filepath.Join(path, "fixtures", "sample.txt")
	destinationFile := filepath.Join(path, "sample.txt")

	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)
		return
	}

	tt := []struct {
		name       string
		method     string
		body       string
		want       string
		statusCode int
	}{
		{
			name:       "without POST body",
			method:     http.MethodPost,
			body:       "",
			want:       "EOF",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "with actual data to be returned and a correct POST body",
			method:     http.MethodPost,
			body:       `{"filename": "sample.txt", "from": "2021-07-06T23:00:00Z", "to": "2001-07-11T04:43:09Z"}`,
			want:       `[{"eventTime":"2001-07-11T13:16:26Z","email":"test@test.com","sessionId":"9a8f2ac8-5a0b-4acd-8fd8-fd9fcbf244c0"},{"eventTime":"2001-07-12T06:01:06Z","email":"test@test.com","sessionId":"2c2f1ea3-fee5-40b5-a085-92da3845fc4b"},{"eventTime":"2001-07-13T03:28:17Z","email":"test@test.com","sessionId":"cf89cd22-3b6b-46a9-aad3-dfdf5b776da3"},{"eventTime":"2001-07-13T18:38:51Z","email":"test@test.com","sessionId":"11d4ef62-b185-4dfb-87da-1dc4a832e2d0"},{"eventTime":"2001-07-14T17:14:40Z","email":"test@test.com","sessionId":"fc5621fa-212b-4750-8606-7dbc21c94f26"}]`,
			statusCode: http.StatusOK,
		},
		{
			name:       "with bad filename",
			method:     http.MethodPost,
			body:       `{"filename": "txt", "from": "2021-07-06T23:00:00Z", "to": "2001-07-11T04:43:09Z"}`,
			want:       "null",
			statusCode: http.StatusOK,
		},
		{
			name:       "with improper date 1",
			method:     http.MethodPost,
			body:       `{"filename": "sample.txt", "from": "2021-07-06T23:00:00Z", "to": "2001-07-11"}`,
			want:       `parsing time "\"2001-07-11\"" as "\"2006-01-02T15:04:05Z07:00\"": cannot parse "\"" as "T"`,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "with improper date 2",
			method:     http.MethodPost,
			body:       `{"filename": "sample.txt", "from": "2021-07-06T23:00:00Z", "to": ""}`,
			want:       `parsing time "\"\"" as "\"2006-01-02T15:04:05Z07:00\"": cannot parse "\"" as "2006"`,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "with improper date 3",
			method:     http.MethodPost,
			body:       `{"filename": "sample.txt", "from": "", "to": ""}`,
			want:       `parsing time "\"\"" as "\"2006-01-02T15:04:05Z07:00\"": cannot parse "\"" as "2006"`,
			statusCode: http.StatusNotFound,
		},
		{
			name:       "with no values in the JSON",
			method:     http.MethodPost,
			body:       `{"filename": "", "from": "", "to": ""}`,
			want:       `parsing time "\"\"" as "\"2006-01-02T15:04:05Z07:00\"": cannot parse "\"" as "2006"`,
			statusCode: http.StatusNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(tc.method, "/", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			handler := handler.Handler{H: handler.GetResults}
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
