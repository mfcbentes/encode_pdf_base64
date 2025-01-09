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

	// Inserir o PDF codificado no banco de dados
	insertQuery := `
        INSERT INTO laudo_paciente_pdf_serial (
            nr_sequencia, dt_atualizacao, nm_usuario, nr_acesso_dicom, nr_seq_laudo, ds_pdf_serial
        ) VALUES (
            laudo_paciente_pdf_serial_seq.nextval, SYSDATE, :1, :2, :3, :4
        )`
	_, err = db.Exec(insertQuery, "manoel.bentes", 20311, 18557, base64PDF)
	if err != nil {
		log.Fatalf("Erro ao inserir no banco: %v", err)
	}

	fmt.Println("PDF inserido com sucesso no banco!")

	// Consulta para recuperar o PDF codificado
	selectQuery := `
        SELECT
            initcap(obter_nome_pf(ap.cd_pessoa_fisica)) AS nm_paciente,
            TRIM(obter_desc_proc_interno(pp.nr_seq_proc_interno)) AS ds_procedimento,
            ap.cd_pessoa_fisica AS protocolo,
            pp.nr_seq_interno AS senha,
            obter_compl_pf(ap.cd_pessoa_fisica, 1, 'DDDCEL') || obter_compl_pf(ap.cd_pessoa_fisica, 1, 'CEL') AS nr_telefone,
            ge.nr_prescricao,
            ap.nr_atendimento,
            lpps.nr_acesso_dicom,
            lpps.DS_PDF_SERIAL
        FROM
            eis_gestao_exames_v ge,
            prescr_procedimento pp,
            prescr_medica pm,
            atendimento_paciente ap,
            laudo_paciente_pdf_serial lpps
        WHERE
            trunc(lpps.dt_atualizacao) = trunc(sysdate)
            AND ge.nr_prescricao = pp.nr_prescricao
            AND ge.nr_seq_proced = pp.nr_sequencia
            AND pm.nr_prescricao = pp.nr_prescricao
            AND pm.nr_atendimento = ap.nr_atendimento
            AND pp.nr_seq_interno = lpps.nr_acesso_dicom
        ORDER BY
            lpps.dt_atualizacao`

	rows, err := db.Query(selectQuery)
	if err != nil {
		log.Fatalf("Erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			nmPaciente     string
			dsProcedimento string
			protocolo      int
			senha          int
			nrTelefone     string
			nrPrescricao   int
			nrAtendimento  int
			nrAcessoDicom  int
			dsPdfSerial    string
		)

		err := rows.Scan(&nmPaciente, &dsProcedimento, &protocolo, &senha, &nrTelefone, &nrPrescricao, &nrAtendimento, &nrAcessoDicom, &dsPdfSerial)
		if err != nil {
			log.Fatalf("Erro ao escanear a linha: %v", err)
		}

		// Decodificar o PDF
		pdfData, err := base64.StdEncoding.DecodeString(dsPdfSerial)
		if err != nil {
			log.Fatalf("Erro ao decodificar o Base64: %v", err)
		}

		// Salvar o PDF no diretório output
		outputFilePath := fmt.Sprintf("output/%d_%d.pdf", protocolo, senha)
		err = os.WriteFile(outputFilePath, pdfData, 0644)
		if err != nil {
			log.Fatalf("Erro ao escrever o arquivo PDF: %v", err)
		}

		fmt.Printf("PDF decodificado e salvo com sucesso em %s\n", outputFilePath)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Erro nas linhas da consulta: %v", err)
	}
}
