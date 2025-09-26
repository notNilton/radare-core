package db

import (
	"context"
	"fmt"
	"log"
	"os" // Pacote para ler variáveis de ambiente
	"radare-core/backend/backend/ent"

	"github.com/joho/godotenv" // Importa a biblioteca para ler o .env
	_ "github.com/lib/pq"
)

// NewClient cria e retorna um novo cliente Ent para PostgreSQL.
func NewClient() *ent.Client {
	// Carrega as variáveis de ambiente do arquivo .env.
	// Se o arquivo .env não existir, a função não retorna erro,
	// permitindo que as variáveis sejam setadas diretamente no ambiente de produção.
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Erro ao carregar o arquivo .env")
	}

	// Lê os detalhes da conexão a partir das variáveis de ambiente.
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Constrói a string de conexão (DSN) usando os valores carregados.
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=public",
		host, port, user, password, dbname)

	// Abrir conexão com o banco de dados PostgreSQL.
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("falha ao abrir conexão com o postgres: %v", err)
	}

	// Executar as migrações automáticas para criar o schema do banco de dados.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("falha ao criar schema: %v", err)
	}

	log.Println("Conexão com o PostgreSQL e schema atualizados com sucesso.")
	return client
}
