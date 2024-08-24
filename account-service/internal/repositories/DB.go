package repositories

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func NewDBConnection() *sql.DB {
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	name := viper.GetString("db.name")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)

	result, err := sql.Open("postgres", connString)
	if err != nil {
		slog.Error("Error creating db connecation", " error", err.Error())
		panic(err)
	}

	return result
}
