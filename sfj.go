package sfj

import (
	"go/format"
	"strings"
	"sync"
)

// Do performs a GET request to each defined route concatenated with the server name.
// lines is a slice containing routes in the format: `/path/to/request` or
// `/path/:parameter1/:parameter2/request parameter1Value parameter2Value`.
// It passes headers in each request and returns a file whose package is a pkg containing structure definitions.
func Do(pkg, server string, lines []string, headerMap map[string]string, insecure, subStruct bool) ([]byte, error) {
	var wg sync.WaitGroup
	server = strings.TrimRight(server, "/")
	lines = deleteEmpty(lines)
	n := len(lines)
	wg.Add(n)
	c := make(chan result, n)
	defer close(c)

	for i := 0; i < n; i++ {
		go requestConverter(server, lines[i], pkg, headerMap, c, &wg, insecure, subStruct)
	}
	wg.Wait()

	var structs []byte
	for i := 0; i < n; i++ {
		if r := <-c; r.err != nil {
			return r.res, r.err
		} else {
			structs = append(structs, r.res...)
		}
	}

	fileContent := string(structs)
	fileContent = strings.Replace(fileContent, "}\npackage "+pkg, "}", -1)

	return format.Source([]byte(fileContent))
}
