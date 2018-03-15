package myHttp

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

	// proxy request
	var err error
	if r.URL, err = url.Parse(p.Addr + r.URL.String()); nil != err {
		log.Println(err)
		return
	}
	r.RequestURI = ""
	client := &http.Client{}
	var resp *http.Response
	fmt.Println(r.URL)
	if resp, err = client.Do(r); nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	// proxy response
	io.Copy(w, resp.Body)
}
