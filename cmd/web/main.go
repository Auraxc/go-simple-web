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
	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	// 使用 log.New() 方法创建一个日志记录器，需要传入三个参数：日志的输出方，日志前缀，附带的信息（日期和时间，使用"|"连接）
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	// 创建一个错误日志记录器，以 stderr 作为输出方，
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()

	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	// 设置一个新的命令行参数 addr，默认值为 ":4000"，同时添加简短的参数说明，在启动时这个参数将存到 addr 变量中
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	// 使用 flag.Parse() 方法来解析命令行参数
	flag.Parse()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	// 创建一个文件服务，指定 "./ui/static/" 为根目录
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// 初始化一个新的 application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	// 用 mux.Handle() 注册文件服务，以 /static/ 开头的路径都将路由到文件服务
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	infoLog.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	errorLog.Fatal(err)
}
