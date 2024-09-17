package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Gileno29/file_loader_golang/models"

	_ "github.com/lib/pq"
)

func conectar() *sql.DB {

	connStr := "user=seu_usuario dbname=seu_banco password=sua_senha host=localhost sslmode=disable"
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

func inserirRegistros(v models.Venda, c *sql.DB) {

}
