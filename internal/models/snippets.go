package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// Write the SQL statement we want to execute.
	// SQL 语句
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use the Query() method on the connection pool to execute our
	// SQL statement. This returns a sql.Rows resultset containing the result of
	// our query.
	// 用 Query() 方法执行语句，返回查询到的所有行
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	// 调用 defer rows.Close() 保证获取 sql.Rows 在 Latest() 返回前正确调用
	// 如果调用报错，也能保证 rows.Close() 被调用，关闭底层的数据库连接
	// 如果数据库连接不能被正确关闭，那么连接池的连接将被耗尽，之后就无法连接数据库了
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	// 初始化一个空的切片来保存 Snippet 结构体
	snippets := []*Snippet{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan()
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		// 用 Scan() 方法从原始数据复制到 Snippet 结构体中
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		// 添加到切片 snippets
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK then return the Snippets slice.
	// 如果一切正常则返回切片 snippets
	return snippets, nil
}
