package main

import (
	"log"
	"net/http"
)

// 未使用的方法声明为 "_"

func home(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello from Snippetbox"))
	if err != nil {
		return
	}
}

// 添加 snippetView 处理方法
// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Display a specific snippet..."))
	if err != nil {
		return
	}
}

// 添加 snippetCreate 处理方法
// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Create a new snippet..."))
	if err != nil {
		return
	}
}

func main() {
	// 注册两个新的 handler 方法，并且指定路由
	// 方法和之前注册 home 路由是一致的：实现路由方法，注册路由
	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
