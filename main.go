package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
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
		for scanner.Scan() {
			// Substituir os espaços por '|'
			modifiedLine := re.ReplaceAllString(scanner.Text(), ";") + "\n"
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
	count := 0

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

	// Exibir o conteúdo do CSV
	for i, record := range records {
		if count == 0 {
			fmt.Print("estou aqui")
			count = count + 1
			continue
		}

		fmt.Printf("Linha %d: %v\n", i+1, record[0])
		break
	}

}

func main() {

	insertIntoDB("Base.txt")

}
