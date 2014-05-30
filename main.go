package main

import (
    "log"
    "fmt"
    "flag"
    "time"
    "strings"
    "net/http"
    "html/template"

    "github.com/jmervine/gojson"
)

var (
    Listen string
    Port int
    Template, err = template.ParseFiles("./index.html")
    defaultJson = `{ "key": "val" }`
)

type Result struct {
    Json, Struct string
}

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    begin := time.Now()

    defer r.Body.Close()

    log.Printf("=> %v %v 200 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(begin)))

    res := Result{
        Json: defaultJson,
    }

    if r.Method == "POST" {
        val := r.PostFormValue("json")
        res.Json = val
    }

    if out, e := json2struct.Generate(strings.NewReader(res.Json), "MyJsonName", "main"); e == nil {
        res.Struct = string(out)
    } else {
        log.Printf("=> %v %v 500 %v %v %s\n", r.Method, r.URL, r.Proto, r.Header["User-Agent"], fmt.Sprintf("%s", time.Since(begin)))
        log.Printf("ERROR:\n%v\n\n", e)
        res.Struct = fmt.Sprintf("JSON Parse Error: %v\n", e)
    }
    Template.Execute(w, res)
}

func main() {
    flag.IntVar(&Port, "port", 8080, "startup port")
    flag.StringVar(&Listen, "listen", "localhost", "listen address")
    flag.Parse()
    if err != nil {
        log.Fatal(err)
    }

    handler := Handler{}

    server := &http.Server{
        Addr:           fmt.Sprintf("%s:%d", Listen, Port),
        Handler:        handler,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    log.Printf("Starting at %v\n", server.Addr)
    log.Fatal(server.ListenAndServe())
}