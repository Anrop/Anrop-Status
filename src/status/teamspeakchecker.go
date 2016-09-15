package status

import (
	"fmt"
	"time"

	"github.com/Darfk/ts3"
	"github.com/sourcegraph/checkup"
)

// TeamspeakChecker implements a Checker for Teamspeak servers.
type TeamspeakChecker struct {
	// Name is the name of the endpoint.
	Name string `json:"endpoint_name"`

	// URL is the URL of the endpoint.
	URL string `json:"endpoint_url"`

	// ThresholdRTT is the maximum round trip time to
	// allow for a healthy endpoint. If non-zero and a
	// request takes longer than ThresholdRTT, the
	// endpoint will be considered unhealthy. Note that
	// this duration includes any in-between network
	// latency.
	ThresholdRTT time.Duration `json:"threshold_rtt,omitempty"`

	// Attempts is how many requests the client will
	// make to the endpoint in a single check.
	Attempts int `json:"attempts,omitempty"`
}

// Check performs checks using c according to its configuration.
// An error is only returned if there is a configuration error.
func (c TeamspeakChecker) Check() (checkup.Result, error) {
	if c.Attempts < 1 {
		c.Attempts = 1
	}

	result := checkup.Result{Title: c.Name, Endpoint: c.URL, Timestamp: checkup.Timestamp()}
	result.Times = c.doChecks()

	return c.conclude(result), nil
}

// doChecks executes req using c.Client and returns each attempt.
func (c TeamspeakChecker) doChecks() checkup.Attempts {
	checks := make(checkup.Attempts, c.Attempts)
	for i := 0; i < c.Attempts; i++ {
		start := time.Now()
		client, err := ts3.NewClient(c.URL)
		if err != nil {
			checks[i].Error = err.Error()
			continue
		}

		_, err = client.Exec(ts3.Version())
		if err != nil {
			checks[i].Error = err.Error()
			continue
		}

		client.Close()

		checks[i].RTT = time.Since(start)
		if err != nil {
			checks[i].Error = err.Error()
			continue
		}
	}

	return checks
}

// conclude takes the data in result from the attempts and
// computes remaining values needed to fill out the result.
// It detects degraded (high-latency) responses and makes
// the conclusion about the result's status.
func (c TeamspeakChecker) conclude(result checkup.Result) checkup.Result {
	result.ThresholdRTT = c.ThresholdRTT

	// Check errors (down)
	for i := range result.Times {
		if result.Times[i].Error != "" {
			result.Down = true
			return result
		}
	}

	// Check round trip time (degraded)
	if c.ThresholdRTT > 0 {
		stats := result.ComputeStats()
		if stats.Median > c.ThresholdRTT {
			result.Notice = fmt.Sprintf("median round trip time exceeded threshold (%s)", c.ThresholdRTT)
			result.Degraded = true
			return result
		}
	}

	result.Healthy = true
	return result
}
