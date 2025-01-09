package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	_ "github.com/godror/godror"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbService := os.Getenv("DB_SERVICE")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbService == "" {
		log.Fatal("Faltam variáveis no arquivo .env")
	}

	dsn := fmt.Sprintf("%s/%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbService)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	// Caminho do arquivo PDF
	pdfPath := "input/laudo.pdf"

	// Leia o arquivo PDF
	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo PDF: %v", err)
	}

	// Converta o conteúdo para Base64
	base64PDF := base64.StdEncoding.EncodeToString(pdfData)

	query := `
		INSERT INTO laudo_paciente_pdf_serial (
			nr_sequencia, dt_atualizacao, nm_usuario, nr_acesso_dicom, nr_seq_laudo, ds_pdf_serial
		) VALUES (
			laudo_paciente_pdf_serial_seq.nextval, SYSDATE, :1, :2, :3, :4
		)`
	_, err = db.Exec(query, "manoel.bentes", 20311, 18557, base64PDF)
	if err != nil {
		log.Fatalf("Erro ao inserir no banco: %v", err)
	}

	fmt.Println("PDF inserido com sucesso no banco!")
}
