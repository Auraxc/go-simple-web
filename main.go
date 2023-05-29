package main

import (
	"log"
	"net/http"
)

// 未使用的方法声明为 "_"

func home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello from Snippetbox"))
	// 如果路由没有匹配到，最终都会路由到 “/” 下
	// 对于不存在的路由需要返回一个 404 页面，在这里判断一下路径，如果不是 “/” 说明是一个不存在的路由
	// 返回一个 404 页面
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Create a new snippet..."))
	if r.Method != "POST" {
		// 这里使用一个判断来限定请求方法，如果不是 POST 请求，则返回一个 405 状态码，并且返回一个“请求方法不允许”的提示。
		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.

		// 如果需要发送一个状态码不是 200 的响应，上一个提交中的方法有一些繁琐，这里可以使用 http.Error() 进行简化
		// http.Error() 会进行一系列操作，基本等价于之前进行的操作
		w.Header().Set("Allow", "POST")

		// 显式地声明状态码
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
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
