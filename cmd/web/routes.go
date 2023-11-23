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

	// 不需要登录验证的路由使用 dynamic 中间件链
	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Update these routes to use the new dynamic middleware chain followed by
	// the appropriate handler function. Note that because the alice ThenFunc()
	// method returns a http.Handler (rather than a http.HandlerFunc) we also
	// need to switch to registering the route using the router.Handler() method.
	// 更新路由来使用新的 dynamic 中间件，因为 ThenFunc() 方法返回一个 http.Handler，我们需要使用 Handler 替代 HandlerFunc
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	// 需要登录验证的路由使用 protected 中间件调用链
	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
