package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

type MysqlClient struct {
	MYSQL *sql.DB
}

func (m Mysql) Connect() *MysqlClient {

	mysql, _ := sql.Open("mysql", m.Username+":"+m.Password+"@tcp("+m.Host+":"+m.Port+")/"+m.DatabaseName+"?parseTime=true")
	client := &MysqlClient{MYSQL: mysql}
	return client

}

func (m *MysqlClient) Disconnect() {
	m.Disconnect()
	log.Println("Db Closed")
}
