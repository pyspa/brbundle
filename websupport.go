package brbundle


import (
	"strings"
	"strconv"
	"fmt"
)

const DefaultQValue = 1.0

// based on https://github.com/NYTimes/gziphandler/blob/master/gzip.go
func supportBrotli(s string) bool {
	var e []string

	for _, ss := range strings.Split(s, ",") {
		coding, qvalue, err := parseCoding(ss)

		if err != nil {
			e = append(e, err.Error())
		} else {
			if coding == "br" {
				return qvalue > 0.0
			}
		}
	}
	return false
}

// parseCoding parses a single conding (content-coding with an optional qvalue),
// as might appear in an Accept-Encoding header. It attempts to forgive minor
// formatting errors.
func parseCoding(s string) (coding string, qvalue float64, err error) {
	for n, part := range strings.Split(s, ";") {
		part = strings.TrimSpace(part)
		qvalue = DefaultQValue

		if n == 0 {
			coding = strings.ToLower(part)
		} else if strings.HasPrefix(part, "q=") {
			qvalue, err = strconv.ParseFloat(strings.TrimPrefix(part, "q="), 64)

			if qvalue < 0.0 {
				qvalue = 0.0
			} else if qvalue > 1.0 {
				qvalue = 1.0
			}
		}
	}

	if coding == "" {
		err = fmt.Errorf("empty content-coding")
	}

	return
}
