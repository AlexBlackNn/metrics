package tests

import (
	"bytes"
	"net/http"
)

func (ms *MetricsSuite) TestServerHappyPathV2() {

	type Want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		body []byte
		want Want
	}{
		{
			name: "gauge with value 10.3",
			url:  "/update/",
			body: []byte(
				`{
				"id": "test_gauge",
				"type": "gauge",
				"value": 10.3
				}`,
			),
			want: Want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
		{
			name: "counter with value 10",
			url:  "/update/",
			body: []byte(
				`{
				"id": "test_counter",
				"type": "counter",
				"delta": 10
				}`,
			),
			want: Want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	// stop server when tests finished
	defer ms.srv.Close()

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			url := ms.srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			ms.Equal(tt.want.code, res.StatusCode)
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func (ms *MetricsSuite) TestServerWrongType() {

	type Want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		url  string
		body []byte
		want Want
	}{
		{
			name: "gauge with value 10.3",
			url:  "/update/",
			body: []byte(
				`{
				"id": "test_gauge",
				"type": "NOTgauge",
				"value": 10.3
				}`,
			),
			want: Want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
		{
			name: "counter with value 10",
			url:  "/update/",
			body: []byte(
				`{
				"id": "test_counter",
				"type": "NOTcounter",
				"delta": 10
				}`,
			),
			want: Want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
			},
		},
	}
	// stop server when tests finished
	defer ms.srv.Close()

	for _, tt := range tests {
		ms.Run(tt.name, func() {
			url := ms.srv.URL + tt.url
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(tt.body))
			ms.NoError(err)
			res, err := ms.client.Do(request)
			ms.NoError(err)
			ms.Equal(tt.want.code, res.StatusCode)
			defer res.Body.Close()
			ms.Equal(tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
