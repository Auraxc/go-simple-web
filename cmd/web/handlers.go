package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// 如果路由没有匹配到，最终都会路由到 “/” 下
	// 对于不存在的路由需要返回一个 404 页面，在这里判断一下路径，如果不是 “/” 说明是一个不存在的路由
	// 返回一个 404 页面
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	// 初始化一个包含两个文件路径的切片，注意：定义的基础模板必须在第一个文件的位置
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and use
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	// 用 template.ParseFiles() 方法来读取模板到 template set，如果报错就返回一个 500 错误
	// 路径必须是基于工作目录的相对路径，或者是绝对路径
	// 注意这里的参数与之前不同
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	// 用 ExecuteTemplate() 方法
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// 添加 snippetView 处理方法
// Add a snippetView handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404 page
	// not found response.

	// 通过 r.URL.Query().Get() 方法获取 id 参数
	// 因为用户输入的内容是不可信的，这里使用 strconv.Atoi() 方法将 id 转换成 int 类型的数字
	// 如果转换失败或者 id 小于1，则返回一个 404 页面
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		return
	}
	// Use the fmt.Fprintf() function to interpolate the id value with our response
	// and write it to the http.ResponseWriter.
	// 使用 fmt.Fprintf() 方法将 id 添加到响应中
	// 注意这里使用了一个全新的返回响应的方法，fmt.Fprintf()，其第一个参数为 io.Writer(), 这个参数是一个接口
	// 因为 http.ResponseWriter() 实现了这个接口，所以将它当作参数传入没有问题
	_, err = fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	if err != nil {
		return
	}
}

// 添加 snippetCreate 处理方法
// Add a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
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

	// 当返回请求时，Go 会自动帮你填充以下响应头：Date、Content-Length、Content-Type
	// Go 会主动猜测响应地类型，如果匹配不到，最终会设置为 Content-Type: application/octet-stream
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte(`{"name":"Alex"}`))

	// 设置响应头的一些方法
	// 请求头会自动规范化为第一个字母大写，其余字母小写的形式，如果不需要这种自动格式化的操作，可以使用如下声明
	//w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	// Set a new cache-control header. If an existing "Cache-Control" header exists
	// it will be overwritten.
	w.Header().Set("Cache-Control", "public, max-age=31536000")

	// In contrast, the Add() method appends a new "Cache-Control" header and can
	// be called multiple times.
	w.Header().Add("Cache-Control", "public")
	w.Header().Add("Cache-Control", "max-age=31536000")

	// Delete all values for the "Cache-Control" header.
	// Del 操作不会删除 Go 自动设置的响应头，如果需要删除，那么需要使用如下方法
	// w.Header()["Date"] = nil
	w.Header().Del("Cache-Control")

	// Retrieve the first value for the "Cache-Control" header.
	w.Header().Get("Cache-Control")

	// Retrieve a slice of all values for the "Cache-Control" header.
	w.Header().Values("Cache-Control")

	_, err := w.Write([]byte("Create a new snippet..."))

	if err != nil {
		return
	}
}
