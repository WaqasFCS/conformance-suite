package results

import (
	"encoding/json"
	"time"

	"bitbucket.org/openbankingteam/conformance-suite/pkg/model"
	"gopkg.in/resty.v1"
)

type Metrics struct {
	TestCase     *model.TestCase
	ResponseTime time.Duration // Http Response Time
	ResponseSize int           // Size in bytes of the HTTP Response body
}

// MarshalJSON is a custom marshaler which formats a Metrics struct
// with a response time represented as milliseconds
// response time precision is up to 6 decimal places of precision (nanosecond)
func (m Metrics) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ResponseTime float64 `json:"response_time"`
		ResponseSize int     `json:"response_size"`
	}{
		ResponseTime: float64(m.ResponseTime) / float64(time.Millisecond),
		ResponseSize: m.ResponseSize,
	})
}

var NoMetrics = Metrics{}

func NewMetricsFromRestyResponse(testCase *model.TestCase, response *resty.Response) Metrics {
	return NewMetrics(testCase, response.Time(), len(response.Body()))
}

func NewMetrics(testCase *model.TestCase, responseTime time.Duration, responseSize int) Metrics {
	return Metrics{
		TestCase:     testCase,
		ResponseTime: responseTime,
		ResponseSize: responseSize,
	}
}
