package main

type User struct {
	ID        int    `db_field:"id" db_type:"SERIAL PRIMARY KEY"`
	FirstName string `db_field:"first_name" db_type:"VARCHAR(100)"`
	LastName  string `db_field:"last_name" db_type:"VARCHAR(100)"`
	Email     string `db_field:"email" db_type:"VARCHAR(100) UNIQUE"`
}

func (u *User) TableName() string {
	return "users"
}

type Tabler interface {
	TableName() string
}
type SQLGenerator interface {
	CreateTableSQL(table Tabler) string
	CreateInsertSQL(model Tabler) string
}

//func NewMigrator(db *sql.DB, sqlGenerator SQLGenerator) *Migrator {
//	return &Migrator{
//		db:           db,
//		sqlGenerator: sqlGenerator,
//	}
//}

// Основная функция
func main() {

	//// Подключение к SQLite БД
	//db, err := sql.Open("sqlite3", "file:my_database.db?cache=shared&mode=rwc")
	//if err != nil {
	//	log.Fatalf("failed to connect to the database: %v", err)
	//}
	//// Создание мигратора с использованием вашего SQLGenerator
	//migrator := NewMigrator(db, YourSQLGeneratorInstance)
	//// Миграция таблицы User
	//if err := migrator.Migrate(&User{}); err != nil {
	//	log.Fatalf("failed to migrate: %v", err)
	//}
}
