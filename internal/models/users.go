package models

import (
	"database/sql"
	"time"
)

// Define a new User type. Notice how the field names and types align
// with the columns in the database "users" table?
// 定义 User 类型，字段和数据库表中的是一致的
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Define a new UserModel type which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// We'll use the Insert method to add a new record to the "users" table.
// 用 Insert 方法添加一个新记录到 users 表
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
// 使用 Authenticate 方法来验证用户是否存在，如果用户存在返回用户 id
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
// 使用 Exists 方法检查特定用户是否存在
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
