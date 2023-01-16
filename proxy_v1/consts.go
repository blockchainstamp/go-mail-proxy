package proxy_v1

import "time"

const (
	StatusCode220 = 220
	StatusCode250 = 250

	DialTimeout       = 30 * time.Second
	CommandLineLimit  = 1 << 14 //RFC 5321 (Section 4.5.3.1.6)
	CommandTimeout    = 5 * time.Minute
	SubmissionTimeout = 12 * time.Minute
)
