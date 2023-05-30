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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// 初始化一个新的 application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// 初始化 http server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
