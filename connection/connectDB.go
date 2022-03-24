package connection

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type dbENV struct {
	host     string
	port     int
	username string
	password string
	dbName   string
	sslmode  string
}

func getENV() dbENV {
	err := godotenv.Load(".env")
	if err != nil {
		panic("env file not set")
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("improper port config")
	}
	// default sslmode to disable
	sslmode := os.Getenv("SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	return dbENV{
		host:     os.Getenv("HOST"),
		port:     port,
		username: os.Getenv("USERNAME"),
		password: os.Getenv("PASSWORD"),
		dbName:   os.Getenv("DBNAME"),
		sslmode:  sslmode,
	}
}

func CreateConnection() {
	config := getENV()
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.host, config.port, config.username, config.password, config.dbName, config.sslmode)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		panic(fmt.Sprintf("Error in db init, %v", err))
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(200)

	DB = db
}
