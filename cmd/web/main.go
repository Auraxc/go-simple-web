package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
// 定义一个 application 结构，这里只先定义两个 log 方法，之后还会添加更多东西
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// 经过简化的 main 功能限制为：1.解析运行所需的参数，2.处理依赖关系，3.运行 HTTP server
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// 初始化一个新的 application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// 初始化 http server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servemux containing our routes.
		// 使用 routes() 初始化路由
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
