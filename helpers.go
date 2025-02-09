package sfj

import (
	"bytes"
	"strings"
)

// replaceParameters replaces parameters in the route string,
// if present returns the path with the replaced parameter and
// the path with the parameter (without ':') to build the name.
func replaceParameters(line string) (string, string) {
	line = strings.TrimSpace(line)
	if !strings.Contains(line, " ") {
		return line, line
	}

	var ret bytes.Buffer
	var parametersFound int
	pathAndParams := strings.Split(line, " ")
	path := pathAndParams[0]
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			for i < l && path[i] != '/' {
				i++
			}

			// replace parameter with its value
			parametersFound++
			ret.WriteString(pathAndParams[parametersFound])
			if i < l && path[i] == '/' {
				ret.WriteRune('/')
			}
		} else {
			ret.WriteByte(path[i])
		}
	}

	return ret.String(), strings.Replace(path, ":", "", -1)
}
