package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Vendas struct {
	cpf                string
	private            int32
	incompleto         int32
	ultimaCompra       string
	ticketMedio        string
	ticketUltimaCompra string
	lojaMaisFrequente  string
	lojaUltimaCompra   string
}

func readFile(fileName string) {
	count := 0
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	// Ler o arquivo CSV com delimitador personalizado
	reader := csv.NewReader(file)
	reader.Comma = ';' // Definindo o delimitador como ponto e vírgula

	// Ler todas as linhas
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Erro ao ler o arquivo CSV:", err)
		return
	}

	// Exibir o conteúdo do CSV
	for i, record := range records {
		if count == 0 {
			continue
		}

		fmt.Printf("Linha %d: %v\n", i+1, record)
	}

}
