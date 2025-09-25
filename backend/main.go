// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"radare-core/backend/internal/db"
	"radare-core/backend/internal/handlers"
	"radare-core/backend/internal/middleware"
	"syscall"
	"time"
)

// main é o ponto de entrada da aplicação.
// Ele inicializa o banco de dados, configura as rotas da API,
// inicia o servidor HTTP e lida com o desligamento gracioso.
func main() {
	// Inicializa o cliente do banco de dados Ent.
	// A função NewClient também lida com as migrações do esquema.
	client := db.NewClient()
	defer client.Close()

	// Cria uma nova instância dos manipuladores da aplicação, injetando o cliente do banco de dados.
	h := handlers.New(client)

	// Registra os manipuladores de rota para os endpoints da API.
	// O ErrorHandler é um middleware que centraliza o tratamento de erros.
	http.HandleFunc("/current-values", middleware.ErrorHandler(h.GetCurrentValues))
	http.HandleFunc("/values/history", middleware.ErrorHandler(h.GetValueHistory))
	http.HandleFunc("/healthz", middleware.ErrorHandler(handlers.HealthCheck))

	// Determina a porta do servidor a partir da variável de ambiente PORT, ou usa 8080 como padrão.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configura o servidor HTTP.
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux, // Usa o multiplexador de requisições padrão.
	}

	// Configura um canal para ouvir os sinais de interrupção do sistema (SIGINT, SIGTERM)
	// para um desligamento gracioso.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Inicia o servidor HTTP em uma goroutine separada para não bloquear a execução principal.
	go func() {
		log.Println("Servidor iniciado na porta " + port + "...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v\n", err)
		}
	}()

	// Inicia a goroutine que atualiza os valores no banco de dados periodicamente.
	go h.StartValueUpdater()

	// Bloqueia a execução até que um sinal de desligamento seja recebido.
	sig := <-sigChan
	log.Printf("Sinal de desligamento recebido: %v, iniciando desligamento gracioso...\n", sig)

	// Cria um contexto com um timeout de 30 segundos para o desligamento.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Tenta desligar o servidor graciosamente.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erro durante o desligamento do servidor: %v\n", err)
	}

	log.Println("Servidor desligado com sucesso.")
}