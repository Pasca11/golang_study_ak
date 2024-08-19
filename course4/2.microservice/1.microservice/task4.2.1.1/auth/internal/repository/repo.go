package repository

import (
	authv1 "auth/proto/gen/auth"
	"database/sql"
	"errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

var ErrNoUser = errors.New("no user")

type AuthRepository interface {
	FindByUsername(username string) (*authv1.User, error)
	CreateUser(user *authv1.User) error
}

type Storage struct {
	DB *sql.DB
}

func NewStorage() (*Storage, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	dataBase := &Storage{}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	connStr := "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " host=" + dbHost + " port=" + dbPort + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	dataBase.DB = db
	err = dataBase.CreateNewUserTable()
	return dataBase, err
}

func (db *Storage) CreateNewUserTable() error {
	newTableString := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		isExist BOOLEAN DEFAULT true
	);`

	_, err := db.DB.Exec(newTableString)
	return err
}

func (db *Storage) FindByUsername(username string) (*authv1.User, error) {
	stmt, err := db.DB.Prepare("SELECT username, password FROM users WHERE username=$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(username)
	res := &authv1.User{}

	err = row.Scan(&res.Username, &res.Password)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db *Storage) CreateUser(user *authv1.User) error {
	stmt, err := db.DB.Prepare("INSERT INTO users (username, password) VALUES ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, user.Password)
	return err
}
