package xHttp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Staticer struct {
	RootPath     string
	ExtraHandler http.HandlerFunc
}

func (s Staticer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if nil != s.ExtraHandler {
		s.ExtraHandler(w, r)
	}
	//Default handler
	fmt.Println(r.Method, r.RequestURI)
	http.FileServer(http.Dir(s.RootPath)).ServeHTTP(w, r)
}

type Proxyer struct {
	Addr string
	Hook http.HandlerFunc
	// default(reserve path): origin://proxy/user/login => target://proxy/user/login
	// example: /proxy/user/login => /usr/login
	PathRule func(string) string
}

func (p Proxyer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if "OPTIONS" == r.Method {
		fmt.Println("OPTIONS")
		return
	}

	// hook
	if nil != p.Hook {
		p.Hook(w, r)
	}

	// proxy request(default)
	endpoint := r.URL.String()
	if nil != p.PathRule {
		endpoint = p.PathRule(endpoint)
	}
	var err error
	if r.URL, err = url.Parse(p.Addr + endpoint); nil != err {
		log.Println(err)
		return
	}
	r.RequestURI = "" //should empty

	//log
	fmt.Println("Proxy ", r.URL.String())

	client := &http.Client{}
	var resp *http.Response
	blue := "\033[0;31m"
	nc := "\033[0m"
	if resp, err = client.Do(r); nil != err {
		log.Printf(blue+"%v"+nc, err)

		// resp == nil
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	// proxy response

	for key, values := range resp.Header {
		// for _, value := range values {
		// 	log.Println(key, value)
		// 	w.Header().Add(key, value)
		// }
		w.Header().Set(key, values[0])
	}
	// w.Header().Set("Content-Length", "5")
	w.WriteHeader(resp.StatusCode)
	// w.Write([]byte("hello"))
	// w.Write([]byte("d"))
	// go as(w, resp.Body)
	io.Copy(w, resp.Body)
}

// func as(w io.Writer, r io.Reader) {
// 	// io.Copy(w, r)
// 	w.Write([]byte("helllllll"))
// }
