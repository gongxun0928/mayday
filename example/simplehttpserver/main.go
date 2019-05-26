package main

import (
	"fmt"
	"net/http"

	"github.com/gongxun0928/mayday"
)

func AAB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello AAB")
}

func ABB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello ABB")
}

func ABC(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello ABC")
}

func MiddleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello,MiddleTest")
}

func main() {
	mux := mayday.New()
	mux.Use(http.HandlerFunc(MiddleTest))
	mux.GET("/route/aab/:name", http.HandlerFunc(AAB))
	mux.Dump()
	fmt.Println("")
	mux.GET("/route/abb/:name", http.HandlerFunc(ABB))
	mux.Dump()
	fmt.Println("")
	mux.GET("/route/abc/*name", http.HandlerFunc(ABC))
	mux.Dump()
	handlers := mux.GetValue("/route/aab/nihao", "GET")
	for _, h := range handlers {
		h.ServeHTTP(nil, nil)
	}
}
