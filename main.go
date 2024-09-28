package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/klassmann/cpfcnpj"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Venda struct {
	cpf                string
	private            int32
	incompleto         int32
	ultimaCompra       string
	ticketMedio        string
	ticketUltimaCompra string
	lojaMaisFrequente  string
	lojaUltimaCompra   string
	cpfValid           bool
	cnpjValid          bool
}

type Column struct {
	Name string
	Type string
}

type Table struct {
	Name    string
	Columns []Column
}

var conection *sql.DB

func processFile(fileName string) (io.Reader, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	// Criar um pipe para leitura e escrita
	reader, writer := io.Pipe()

	// Usar uma expressão regular para encontrar um ou mais espaços
	re := regexp.MustCompile(`\s+`)

	// Executar a modificação dos dados em uma goroutine
	go func() {
		defer file.Close()
		defer writer.Close()

		scanner := bufio.NewScanner(file)
		count := 0
		for scanner.Scan() {
			if count == 0 {
				count += 1
				continue
			}
			// Substituir os espaços por '|'
			modifiedLine := re.ReplaceAllString(scanner.Text(), "|") + "\n"
			// Escrever a linha modificada no writer
			_, err := writer.Write([]byte(modifiedLine))
			if err != nil {
				fmt.Println("Erro ao escrever no pipe:", err)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
		}
	}()

	// Retornar o lado leitor do pipe
	return reader, nil
}

func insertIntoDB(fileName string, db *sql.DB) {
	//count := 0
	var cpfValid bool
	var cnpjValid bool
	csvMemory, err := processFile(fileName)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo")
	}
	reader := csv.NewReader(csvMemory)
	reader.Comma = '|' // Definindo o delimitador como ponto e vírgula

	// Ler todas as linhas
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo CSV:", err)
		return
	}

	for i, record := range records {

		fmt.Printf("Linha %d: %v\n", i+1, record[2])
		v1, err := strconv.Atoi(record[2])
		if err != nil {
			fmt.Println("Erro ao converter")
		}

		v2, err := strconv.Atoi(record[1])

		if err != nil {
			fmt.Print("Erro o converter valor 02")
		}

		cpf := cpfcnpj.NewCPF(record[0])
		if cpf.IsValid() {
			cpfValid = true

		} else {
			cpfValid = false
		}

		cnpj := cpfcnpj.NewCNPJ(record[7])

		if cnpj.IsValid() {
			cnpjValid = true
		} else {
			cnpjValid = false
		}

		venda := Venda{record[0], int32(v2), int32(v1), record[3], record[4], record[5], record[6], record[7], cpfValid, cnpjValid}

		inserirRegistros(venda, db)
	}

}

func conectar() (error *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	connStr := "user=" + os.Getenv("POSTGRES_USER") + " dbname=" + os.Getenv("POSTGRES_DB") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " host=" + os.Getenv("DATABASE_HOST") + " ?sslmode=disable"
	fmt.Println("String de conexao: ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir a conexão:", err)
	}
	//defer db.Close()

	// Testar a conexão
	/*	err = db.Ping()
		if err != nil {
			log.Fatal("Erro ao conectar ao banco de dados:", err)
		}
	*/
	fmt.Println("Conectado ao banco de dados com sucesso!")

	return db

}

func inserirRegistros(v Venda, c *sql.DB) {

	_, err := c.Exec("INSERT INTO venda VALUES ($1, $2, $3,$4,$5, $6, $7, $8, $9, $10)", v.cpf, v.private, v.incompleto, v.ultimaCompra, v.ticketMedio, v.ticketUltimaCompra, v.lojaMaisFrequente, v.lojaUltimaCompra, v.cpfValid, v.cnpjValid)
	if err != nil {
		log.Fatal("Erro ao executar o INSERT:", err)
	}
	fmt.Println("Registro inserido com sucesso!")

}

func createTable(table Table, c *sql.DB) error {
	var columns []string
	for _, column := range table.Columns {
		columns = append(columns, fmt.Sprintf("%s %s", column.Name, column.Type))
	}

	drop := fmt.Sprintf("DROP TABLE %s ", table.Name)

	_, err := c.Exec(drop)
	if err != nil {
		fmt.Println("Erro ao dropar tabela")
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s)", table.Name, strings.Join(columns, ", "))

	_, err = c.Exec(query)

	if err != nil {
		return err
	}

	return nil

}

//GIN FUNCTIONS

func defaultRouter(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "{message: rota default}")
}

func uploadFile(c *gin.Context) {

	file, err := c.FormFile("arquivo")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ao recuperar arquivo": err.Error()})
		return

	}

	dst, err := os.OpenFile("./uploads/"+file.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ao criar arquivo": err.Error()})
		return
	}
	defer dst.Close()

	openFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if _, err := io.Copy(dst, openFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Arquivo salvo com sucesso"})

	insertIntoDB(dst.Name(), conection)

}
func main() {

	//insertIntoDB("Base.txt")

	//TESTE CONEXAO
	conection = conectar()
	var colunms []Column

	c1 := Column{"cpf", "varchar"}
	c2 := Column{"private", "integer"}
	c3 := Column{"incompleto", "integer"}
	c4 := Column{"ultimaCompra", "varchar"}
	c5 := Column{"ticketmedio", "varchar"}
	c6 := Column{"ticketultimacompra", "varchar"}
	c7 := Column{"lojamaisfrequente", "varchar"}
	c8 := Column{"lojaUltimaCompra", "varchar"}
	c9 := Column{"cpfValid", "bool"}
	c10 := Column{"cnpjValid", "bool"}
	colunms = append(colunms, c1, c2, c3, c4, c5, c6, c7, c8, c9, c10)
	tb := Table{"venda", colunms}

	err := createTable(tb, conection)

	if err != nil {
		fmt.Println(err)
	}

	router := gin.Default()

	router.GET("/", defaultRouter)
	router.POST("/upload", uploadFile)

	router.Run("0.0.0.0:8080")

}
