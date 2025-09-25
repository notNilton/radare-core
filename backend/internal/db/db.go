package db

import (
	"context"
	"log"
	"radare.com/backend/ent"

	_ "github.com/mattn/go-sqlite3"
)

// NewClient cria e retorna um novo cliente Ent.
func NewClient() *ent.Client {
	// Abrir conexão com o banco de dados SQLite.
	// O arquivo do banco de dados será criado se não existir.
	client, err := ent.Open("sqlite3", "file:radare.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("falha ao abrir conexão com o sqlite: %v", err)
	}

	// Executar as migrações automáticas para criar o schema do banco de dados.
	// Não use isso em um ambiente de produção com múltiplos nós.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("falha ao criar schema: %v", err)
	}

	log.Println("Conexão com o banco de dados e schema atualizados com sucesso.")
	return client
}