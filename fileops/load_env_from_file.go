package fileops

import (
	"bufio"
	"os"
	"strings"

	"github.com/rohanthewiz/serr"
)

// EnvFromFile reads a `*.env` style file and loads into the environment
func EnvFromFile(filespec string) (issues []serr.SErr, err error) {
	file, err := os.Open(filespec)
	if err != nil {
		return issues, serr.Wrap(err, "Error reading: "+filespec)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() { // splits on lines by default
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") { // skip lines starting with a comment
			continue
		}

		// Keys and Values
		bef, aft, fnd := strings.Cut(line, "=")
		if !fnd {
			continue
		}

		key := strings.TrimSpace(bef)
		if key == "" {
			issues = append(issues, serr.NewSErr("key is empty", "line", line))
			continue
		}

		val := strings.TrimSpace(aft)
		if val == "" {
			issues = append(issues, serr.NewSErr("Value is empty", "line", line))
			continue
		}

		// Check for delimiters and comments
		if len(val) > 1 {
			// First check if value has surrounding quotes as **quotes have the highest precedence**
			// Don't trim after delimiters removed to allow spaces in values
			if strings.HasPrefix(val, `'`) {
				if idx := strings.IndexByte(val[1:], '\''); idx != -1 {
					val = val[1 : idx+1]
				}
			} else if strings.HasPrefix(val, `"`) {
				if idx := strings.IndexByte(val[1:], '"'); idx != -1 {
					val = val[1 : idx+1]
				}
				// For comments we do want to trim space
			} else if x := strings.IndexByte(val, '#'); x != -1 {
				val = strings.TrimSpace(val[:x])
			}
		}

		if val == "" {
			issues = append(issues, serr.NewSErr("Value is empty", "line", line))
			continue
		}

		os.Setenv(key, val)
	}

	if err := scanner.Err(); err != nil {
		return issues, serr.Wrap(err, "Error while scanning ", filespec)
	}

	return
}
