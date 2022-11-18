package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "021020"
	dbname   = "chatroom"
)

var dbConn *pgx.Conn

func ConnDB() *pgxpool.Pool {
	var err error
	dbConnParam := fmt.Sprintf(`%s://%s:%s@%s:%d/%s`, user, user, password, host, port, dbname)
	pool, err = pgxpool.Connect(context.Background(), dbConnParam)
	// 获取连接池里面的一个连接： conn, err = pool.Acquire(context.Background())
	// 关闭拿出来的那个连接： defer conn.Release()
	if err != nil {
		fmt.Println("failed to connect database")
		//panic(err.Error())
		return nil
	} else {
		fmt.Println("connect database successfully")
	}
	return pool
}

var pool *pgxpool.Pool

func GetConn() *pgxpool.Conn {
	dbConn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	return dbConn
}
