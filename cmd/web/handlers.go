package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method against
// *application.

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	if err != nil {
		return
	}
}

// 添加 snippetCreate 处理方法
// Add a snippetCreate handler function.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		w.Header().Set("Allow", "POST")

		// 显式地声明状态码
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
