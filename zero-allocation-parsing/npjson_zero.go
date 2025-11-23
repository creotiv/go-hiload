package zeroallocationparsing

import (
	"strconv"
)

type LogRecord struct {
	TS  int64
	Msg string
	Lev string
	App string
}

// ParseBatch parses a JSON array of objects with a known schema.
// No reflection, no maps, no allocation except final strings.
func ParseBatchCustom(src []byte) ([]LogRecord, error) {
	var out []LogRecord

	i := 0
	n := len(src)

	// skip whitespace until '['
	for i < n && src[i] != '[' {
		i++
	}
	if i == n {
		return out[:0], nil
	}
	i++ // skip '['

	recordCount := 0

	// main loop: find objects {...}
	for i < n {
		// find '{'
		for i < n && src[i] != '{' {
			if src[i] == ']' {
				return out[:recordCount], nil
			}
			i++
		}
		if i == n {
			break
		}
		i++ // skip '{'

		// ensure out has capacity
		if recordCount >= len(out) {
			out = append(out, LogRecord{})
		}
		rec := &out[recordCount]

		// parse fields
		for i < n && src[i] != '}' {
			// find '"'
			for i < n && src[i] != '"' {
				i++
			}
			i++ // skip first quote

			// key start
			keyStart := i
			for i < n && src[i] != '"' {
				i++
			}
			keyEnd := i
			i++ // skip closing quote

			// skip colon
			for i < n && src[i] != ':' {
				i++
			}
			i++ // skip ':'

			key := string(src[keyStart:keyEnd])

			// parse value depending on type
			switch key {
			case "ts":
				// integer
				for i < n && (src[i] == ' ' || src[i] == '\n') {
					i++
				}
				valStart := i
				for i < n && (src[i] >= '0' && src[i] <= '9') {
					i++
				}
				v, _ := strconv.ParseInt(string(src[valStart:i]), 10, 64)
				rec.TS = v

			case "msg", "lev", "app":
				// string: skip until opening quote
				for i < n && src[i] != '"' {
					i++
				}
				i++ // skip first "
				valStart := i
				for i < n && src[i] != '"' {
					i++
				}
				valEnd := i
				i++ // skip closing "

				s := string(src[valStart:valEnd])

				switch key {
				case "msg":
					rec.Msg = s
				case "lev":
					rec.Lev = s
				case "app":
					rec.App = s
				}
			}

			// skip until next key or end of object
			for i < n && src[i] != '"' && src[i] != '}' {
				i++
			}
		}

		// skip '}'
		if i < n && src[i] == '}' {
			i++
		}
		recordCount++

		// skip comma / whitespace
		for i < n && src[i] != '{' && src[i] != ']' {
			i++
		}
	}

	return out[:recordCount], nil
}
