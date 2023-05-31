package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // New import
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
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string.
	// 这里使用了 docker 中的 mysql，重新映射了端口，所以与教程中的写法不同，注释是原始写法
	//     dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	dsn := flag.String("dsn", "web:pass@tcp(127.0.0.1:13306)/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	// 为了保持 main() 函数的简洁，创建连接池的工作交给 OpenDB() 方法
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	// 延迟调用 db.Close() 方法，连接池可以在 main() 方法结束之前关闭
	// 这个方法在这里永远不会被调用，因为只有 ctrl+c 或 errorLog.Fatal(err) 被触发时 main() 才会退出
	// 但主动关闭建立的连接池是一个好习惯
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := 'declare and assign'
	// operator.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
