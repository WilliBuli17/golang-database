package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO customer(id, name) VALUES ('buli', 'Buli')"
	_, err := db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Intert New Customer")
}

// ---------------------------------------------------------------------------------------------------------------------
func TestQueryContextQsl(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string

		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("id :", id)
		fmt.Println("name :", name)
	}

	defer rows.Close()

	fmt.Println("Success Select New Customer")
}

// ---------------------------------------------------------------------------------------------------------------------
func TestQueryContextQslComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance sql.NullInt32
		var raating sql.NullFloat64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married sql.NullBool

		// tipe data sql.* akan mengembalikan datanya menjadi {value valid} -- value nya adalah nilai dari db, validnya adalah vefifikasi jika datanya itu valid dari db atau tidak, gunakan pengecekan jika hanya ingin ambil valuenya

		err := rows.Scan(&id, &name, &email, &balance, &raating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("-----------------------")
		fmt.Println("id :", id)
		fmt.Println("name :", name)
		if email.Valid {
			fmt.Println("email :", email.String)
		}
		if balance.Valid {
			fmt.Println("balance :", balance.Int32)
		}
		if raating.Valid {
			fmt.Println("raating :", raating.Float64)
		}
		if birthDate.Valid {
			fmt.Println("birth date :", birthDate.Time)
		}
		if married.Valid {
			fmt.Println("married :", married.Bool)
		}
		fmt.Println("created at :", createdAt)
	}

	defer rows.Close()

	fmt.Println("Success Select New Customer")
}

// ---------------------------------------------------------------------------------------------------------------------
func TestSqlInjectionSafeWithParam(t *testing.T) { // query menggunakan parameter, menghindari sql injection
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "admin"
	password := "admin"

	query := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string

		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}

		fmt.Println("username :", username)
	}

	defer rows.Close()

	fmt.Println("Success Select New User")
}

func TestExecSqlInkectionSafeWithParam(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	username := "willi"
	password := "willi"

	query := "INSERT INTO user(username, password) VALUES (?, ?)"
	_, err := db.ExecContext(ctx, query, username, password)

	if err != nil {
		panic(err)
	}

	fmt.Println("Success Intert New User")
}

// ---------------------------------------------------------------------------------------------------------------------
func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	email := "william buli"
	comment := "william buli"

	query := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	result, err := db.ExecContext(ctx, query, email, comment)

	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success Intert New Comments With id :", insertId)
}

// ---------------------------------------------------------------------------------------------------------------------
func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	for i := 0; i < 10; i++ {
		email := "willi" + strconv.Itoa(i) + "@email.xyz"
		comment := "test comment " + strconv.Itoa(i)

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		insertId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Success Intert New Comments With id :", insertId)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO comments(email, comment) VALUES (?, ?)"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		email := "willibuli" + strconv.Itoa(i) + "@email.xyz"
		comment := "test comment tx " + strconv.Itoa(i)

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		insertId, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		fmt.Println("Success Intert New Comments With id :", insertId)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}
