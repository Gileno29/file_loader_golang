package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/Gileno29/file_loader_golang/database"
)

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

func insertIntoDB(fileName string) {
	//count := 0

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
		break
	}

}

func main() {

	//insertIntoDB("Base.txt")

	//TESTE CONEXAO
	c := database.Conectar("uservendas", "vendas", "dbvendas")

	fmt.Println(c)

}
