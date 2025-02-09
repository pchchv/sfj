# SFJ [![GoDoc](https://godoc.org/github.com/pchchv/sfj?status.svg)](https://godoc.org/github.com/pchchv/sfj)

### Struct from JSON

Generator of Go structs from JSON server responses.

SFJ defines type names by using the specified lines in the route file and skipping numbers.
i. e. a request to a route like `/users/1/posts` generates `type UsersPosts`.

SFJ supports *parameters*: a string like `/users/:user/posts/:pid 1 200` generates `type UsersUserPostsPid` from the response to the request `GET /users/1/posts/200`.

SFJ also supports headers personalization, so it can be used to generate types from responses protected by some authorization method.

### Install

`go get -u github.com/pchchv/sfj`

### Usage

```
sfj [options]
  -headers string
    	Headers to add in every request
  -help
    	prints this help
  -insecure
    	Disables TLS Certificate check for HTTPS, use in case HTTPS Server Certificate is signed by an unknown authority
  -out string
    	Output file. Stdout is used if not specified
  -pkg string
    	Package name (default "main")
  -routes string
    	Routes to request. One per line (default "routes.txt")
  -server string
    	sets the server address (default "http://localhost:9090")
  -substruct
    	Creates types for sub-structs
```

### Examples

You can invoke `sfj` by passing a single JSON (anonymous) from stdin and get it converted to a go structure.

```
echo '  {
    "Book Id": 30558257,
    "Title": "Unsouled (Cradle, #1)",
    "Author": "Will Wight",
    "Author l-f": "Wight, Will",
    "Additional Authors": "",
    "BCID": ""
  }' | ./sfj
```

#### obtaining

```go
package main

type Foo1 struct {
        Additional_Authors string `json:"Additional Authors"`
        Author             string `json:"Author"`
        Author_l_f         string `json:"Author l-f"`
        Bcid               string `json:"BCID"`
        Book_Id            int64  `json:"Book Id"`
        Title              string `json:"Title"`
}
```

Or you can set up a more complex scenario by defining a `routes.txt` file with a line for each (parametric) request, and use it as shown below.

*routes.txt*:

```txt
/
/repos/:user/:repo pchchv sfj
```

Run:

```sh
sfj -server https://api.github.com -pkg example
```

#### [Returns](./expected_out.txt)