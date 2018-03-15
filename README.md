# wangrongda/xhttp

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/wangrongda/xhttp
```

## Examples

Start a web server with path '/' for static file, and '/proxy' for proxy:

```go
func main() {
	rootPath := flag.String("p", os.Getenv("PWD"), "Static file server path")
	http.Handle("/", xHttp.Staticer{RootPath: *rootPath})
	http.Handle("/proxy", xHttp.Proxyer{Addr: "http://127.0.0.1:9090", Hook: http.HandlerFunc(handle)})
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("auth", "9234290134813")
}
```

The Staticer struct has two field:
* RootPath: Specify path
* ExtraHandler: If this is not nil, the function will be called before responsing a static source.

```go
type Staticer struct {
	RootPath     string
	ExtraHandler http.HandlerFunc
}
```

The Proxyer struct has two field: 
* Addr: Proxy target ip
* Hook: Function called before requesting target, such as add request head, modify request body...

```go
type Proxyer struct {
	Addr string
	Hook http.HandlerFunc
}
```
