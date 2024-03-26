package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type myHandler struct{}

// GO WORLD => interfaces are implicitly implemented, type has all the methods specified by an interface => it's considered to implement that interface.
// In this case, myHandler is said to implement the http.Handler interface because it has the ServeHTTP method with the correct signature.
//
// from 'server.go'
// type Handler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
