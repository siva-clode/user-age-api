package dbinit

import "database/sql"

func EnsureTables(conn *sql.DB) error {
	createUserTable := `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    dob DATE NOT NULL
);`
	_, err := conn.Exec(createUserTable)
	return err
}
