package healthcheck

import (
	"context"
	"time"

	"golang.org/x/exp/slices"

	"github.com/coder/coder/v2/coderd/database"
	"github.com/coder/coder/v2/coderd/healthcheck/health"
)

const (
	DatabaseDefaultThreshold = 15 * time.Millisecond
)

// @typescript-generate DatabaseReport
type DatabaseReport struct {
	// Healthy is deprecated and left for backward compatibility purposes, use `Severity` instead.
	Healthy   bool             `json:"healthy"`
	Severity  health.Severity  `json:"severity" enums:"ok,warning,error"`
	Warnings  []health.Message `json:"warnings"`
	Dismissed bool             `json:"dismissed"`

	Reachable   bool    `json:"reachable"`
	Latency     string  `json:"latency"`
	LatencyMS   int64   `json:"latency_ms"`
	ThresholdMS int64   `json:"threshold_ms"`
	Error       *string `json:"error"`
}

type DatabaseReportOptions struct {
	DB        database.Store
	Threshold time.Duration

	Dismissed bool
}

func (r *DatabaseReport) Run(ctx context.Context, opts *DatabaseReportOptions) {
	r.Warnings = []health.Message{}
	r.Severity = health.SeverityOK
	r.Dismissed = opts.Dismissed

	r.ThresholdMS = opts.Threshold.Milliseconds()
	if r.ThresholdMS == 0 {
		r.ThresholdMS = DatabaseDefaultThreshold.Milliseconds()
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pingCount := 5
	pings := make([]time.Duration, 0, pingCount)
	// Ping 5 times and average the latency.
	for i := 0; i < pingCount; i++ {
		pong, err := opts.DB.Ping(ctx)
		if err != nil {
			r.Error = health.Errorf(health.CodeDatabasePingFailed, "ping database: %s", err)
			r.Severity = health.SeverityError

			return
		}
		pings = append(pings, pong)
	}
	slices.Sort(pings)

	// Take the median ping.
	latency := pings[pingCount/2]
	r.Latency = latency.String()
	r.LatencyMS = latency.Milliseconds()
	if r.LatencyMS >= r.ThresholdMS {
		r.Severity = health.SeverityWarning
		r.Warnings = append(r.Warnings, health.Messagef(health.CodeDatabasePingSlow, "median database ping above threshold"))
	}
	r.Healthy = true
	r.Reachable = true
}
