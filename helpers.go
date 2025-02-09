package sfj

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/ChimeraCoder/gojson"
)

var (
	unnamedStruct int
	mutex         sync.Mutex
)

type result struct {
	res []byte
	err error
}

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

func requestConverter(server, line, pkg string, headerMap map[string]string, c chan result, wg *sync.WaitGroup, insecure, subStruct bool) {
	// decrement the counter when goroutine ends
	defer wg.Done()

	requestPath, parametricRequest := replaceParameters(line)
	req, _ := http.NewRequest("GET", server+requestPath, nil)
	// set headers
	for key, value := range headerMap {
		req.Header.Set(key, value)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	var err error
	var res *http.Response
	client := &http.Client{Transport: tr}
	if res, err = client.Do(req); err != nil {
		c <- result{nil, err}
		return
	}
	// close writer on end
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c <- result{
			nil,
			fmt.Errorf("Request: %s%s returned status %d\n", server, requestPath, res.StatusCode),
		}
		return
	}

	// generate structName
	var structName string
	requestPathElements := strings.Split(parametricRequest, "/")
	for _, element := range requestPathElements {
		if _, err = strconv.ParseInt(element, 10, 64); err == nil {
			// skip
			continue
		}

		var firstFound bool
		for _, r := range element {
			if unicode.IsLetter(r) {
				if firstFound {
					structName += string(r)
				} else {
					structName += string(unicode.ToUpper(r))
					firstFound = true
				}
			}
		}
	}

	if structName == "" {
		mutex.Lock()
		unnamedStruct++
		structName = "Foo" + strconv.Itoa(unnamedStruct)
		mutex.Unlock()
	}

	var r result
	tagList := []string{"json"}
	convertFloats := true
	r.res, r.err = gojson.Generate(res.Body, gojson.ParseJson, structName, pkg, tagList, subStruct, convertFloats)
	c <- r
}

func deleteEmpty(strs []string) (res []string) {
	for _, str := range strs {
		if str != "" {
			res = append(res, str)
		}
	}
	return
}

func jsonToStruct(pkg, json string) (res result) {
	mutex.Lock()
	unnamedStruct++
	structName := "Foo" + strconv.Itoa(unnamedStruct)
	mutex.Unlock()

	tagList := []string{"json"}
	subStruct, convertFloats := true, true
	res.res, res.err = gojson.Generate(strings.NewReader(json), gojson.ParseJson, structName, pkg, tagList, subStruct, convertFloats)
	return
}
