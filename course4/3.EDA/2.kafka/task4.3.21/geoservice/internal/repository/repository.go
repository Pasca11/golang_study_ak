package repository

import (
	"context"
	"database/sql"
	"fmt"
	models "geoservice-rabbit/models"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresDataBase struct {
	Base *sql.DB
}

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByName(ctx context.Context, name string) (models.User, bool, error)
	Update(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.User, error)
}

func StartPostgresDataBase(ctx context.Context) (*PostgresDataBase, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	dataBase := &PostgresDataBase{}
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

	dataBase.Base = db
	err = dataBase.CreateNewUserTable()
	return dataBase, err
}

func (db *PostgresDataBase) CreateNewUserTable() error {
	newTableString := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		isExist BOOLEAN DEFAULT true
	);`

	_, err := db.Base.Exec(newTableString)
	return err
}

func (db *PostgresDataBase) Create(ctx context.Context, user models.User) error {
	query := `
        INSERT INTO users (username, password, isexist)
        VALUES ($1, $2, $3)
        ON CONFLICT (username) DO NOTHING;
    `

	result, err := db.Base.Exec(query, user.Username, user.Password, true)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with username %s already exists", user.Username)
	}

	return nil
}

func (db *PostgresDataBase) GetByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	query := `SELECT username, password FROM users WHERE id = $1`

	row := db.Base.QueryRow(query, id)
	err := row.Scan(&user.Username, &user.Password)

	if err != nil {

		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with ID %s not found", id)
		}

		return user, err
	}

	return user, nil
}

func (db *PostgresDataBase) Update(ctx context.Context, user models.User) error {
	query := `UPDATE users SET isExist = $1 WHERE username = $2`

	_, err := db.Base.Exec(query, true, user.Username)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresDataBase) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET isExist = $1 WHERE id = $2`

	_, err := db.Base.Exec(query, false, id)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresDataBase) List(ctx context.Context) ([]models.User, error) {
	query := `SELECT username, password, isExist FROM users`
	rows, err := db.Base.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var isExist bool
		err := rows.Scan(&user.Username, &user.Password, &isExist)

		if err != nil {
			return nil, err
		}

		if isExist {
			users = append(users, user)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *PostgresDataBase) GetByName(ctx context.Context, name string) (models.User, bool, error) {
	query := `SELECT username, password, isExist FROM users WHERE username = $1`
	var user models.User

	var isExist bool
	row := db.Base.QueryRow(query, name)
	err := row.Scan(&user.Username, &user.Password, &isExist)

	if err != nil {
		if err == sql.ErrNoRows || !isExist {
			return user, false, fmt.Errorf("user %s not found", name)
		}
	}

	return user, isExist, err
}
