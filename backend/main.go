// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"radare.com/backend/internal/db"
	"radare.com/backend/internal/handlers"
	"radare.com/backend/internal/middleware"
	"syscall"
	"time"
)

func main() {
	// Inicializa o cliente do banco de dados
	client := db.NewClient()
	defer client.Close()

	// Cria uma instância dos handlers, injetando o cliente do banco de dados
	h := handlers.New(client)

	// Define os manipuladores para as rotas, usando o middleware para tratamento de erros
	http.HandleFunc("/api/current-values", middleware.ErrorHandler(h.GetCurrentValues))
	http.HandleFunc("/api/values/history", middleware.ErrorHandler(h.GetValueHistory))
	http.HandleFunc("/healthz", middleware.ErrorHandler(handlers.HealthCheck)) // HealthCheck pode permanecer sem estado

	// Obtém a porta da variável de ambiente PORT ou usa 8080 como padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Cria um servidor HTTP
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
	}

	// Canal para receber sinais de interrupção ou terminação
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Inicia o servidor em uma goroutine
	go func() {
		log.Println("Servidor iniciado na porta " + port + "...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v\n", err)
		}
	}()

	// Inicia a goroutine para atualizar os valores no banco de dados
	go h.StartValueUpdater()

	// Aguarda o sinal de interrupção
	sig := <-sigChan
	log.Printf("Sinal de desligamento recebido: %v, iniciando graceful shutdown...\n", sig)

	// Cria um contexto com timeout para o shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Inicia o processo de shutdown do servidor
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erro durante o shutdown do servidor: %v\n", err)
	}

	log.Println("Servidor desligado com sucesso.")
}