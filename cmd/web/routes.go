package main

import (
	"github.com/julienschmidt/httprouter" // New import
	"github.com/justinas/alice"

	"net/http"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
// 修改 routes() 的签名，返回 http.Handler 替代 *http.ServeMux。
func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Create a handler function which wraps our notFound() helper, and then
	// assign it as the custom handler for 404 Not Found responses. You can also
	// set a custom handler for 405 Method Not Allowed responses by setting
	// router.MethodNotAllowed in the same way too.
	// 使用 app.notFound(w) 来统一 404 页面的相应
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	// 创建一个中间件的调用链
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
