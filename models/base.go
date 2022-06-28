package models

import (
	"database/sql"
	"fmt"

	_ "github.com/nakagami/firebirdsql"
)

var sql_path = "SYSDBA:masterkey@192.168.30.237/d:/database/rehabase.fdb"

func CreateAccount(account *Account) {
	var n int
	conn, _ := sql.Open("firebirdsql", sql_path)
	defer conn.Close()
	conn.QueryRow("SELECT count(*) FROM accounts where email='" + account.Email + "'").Scan(&n)
	if n > 0 {
		account.ID = 0
	} else {
		rows, _ := conn.Query("select GEN_ID(GEN_ACCOUNTS,1) from RDB$DATABASE")
		rows.Next()
		rows.Scan(&account.ID)
		tx, err := conn.Begin()
		if err != nil {
			fmt.Println("Error begin transaction")
		}
		stmt, err := tx.Prepare("INSERT into ACCOUNTS (id,email, passwd) values (?,?,?)")
		if err != nil {
			fmt.Println("Error prepare sql")
		}
		n, err := stmt.Exec(account.ID, account.Email, account.Password)
		if err != nil {
			fmt.Println("Error exec sql", err)
		}
		fmt.Println(n.RowsAffected())
		tx.Commit()
	}
}

func GetAccount(account *Account, email string) bool {
	conn, _ := sql.Open("firebirdsql", sql_path)
	defer conn.Close()
	rows, _ := conn.Query("SELECT ID, email, passwd FROM accounts where email='" + email + "'")
	for rows.Next() {
		rows.Scan(&account.ID, &account.Email, &account.Password)
	}
	if account.Email == email {
		return true
	} else {
		return false
	}

}

func ExistEmail(email string) bool {
	conn, _ := sql.Open("firebirdsql", sql_path)
	defer conn.Close()
	rows, err := conn.Query("SELECT e_mail FROM cards where e_mail='" + email + "'")
	if err == nil {
		email2 := ""
		for rows.Next() {
			rows.Scan(&email2)
			if email2 == email {
				return true
			}
		}
	}
	return false
}

func ValidateUser(email string, password string) bool {
	conn, _ := sql.Open("firebirdsql", sql_path)
	defer conn.Close()
	rows, err := conn.Query("SELECT e_mail, cardcode FROM cards where e_mail='" + email + "'") // 16046
	if err == nil {
		email2 := ""
		password2 := ""
		for rows.Next() {
			rows.Scan(&email2, &password2)
			if password2 == password {
				return true
			}
		}
	}
	return false
}
