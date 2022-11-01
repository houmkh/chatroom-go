package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "021020"
	dbname   = "chatroom"
)

var dbConn *pgx.Conn

func ConnDB() *pgx.Conn {
	var err error
	dbConnParam := fmt.Sprintf(`%s://%s:%s@%s:%d/%s`, user, user, password, host, port, dbname)
	dbConn, err = pgx.Connect(context.Background(), dbConnParam)
	if err != nil {
		fmt.Println("failed to connect database")
		//panic(err.Error())
		return nil
	} else {
		fmt.Println("connect database successfully")
	}
	return dbConn
}
