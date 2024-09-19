package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Gileno29/file_loader_golang/models"

	_ "github.com/lib/pq"
)

func Conectar(user string, pass string, database string) *sql.DB {

	connStr := "user=" + user + " dbname=" + database + " password=" + pass + " host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir a conexão:", err)
	}
	defer db.Close()

	// Testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso!")

	return db

}

func InserirRegistros(v models.Venda, c *sql.DB) {

	_, err := c.Exec("INSERT INTO sua_tabela (nome) VALUES ($1)", v)
	if err != nil {
		log.Fatal("Erro ao executar o INSERT:", err)
	}
	fmt.Println("Registro inserido com sucesso!")

}
