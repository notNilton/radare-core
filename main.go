// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"radare-core/internal/handlers"
	"radare-core/internal/middleware"
	"syscall"
	"time"
)

func main() {
	// Define os manipuladores para as rotas, usando o middleware para tratamento de erros
	http.HandleFunc("/api/current-values", middleware.ErrorHandler(handlers.GetCurrentValues))
	http.HandleFunc("/healthz", middleware.ErrorHandler(handlers.HealthCheck))

	// Obtém a porta da variável de ambiente PORT ou usa 8080 como padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Cria um servidor HTTP
	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux, // Usando o manipulador padrão para as rotas definidas com http.HandleFunc
		// Optionally configure timeouts here, e.g.,
		// ReadTimeout:  15 * time.Second,
		// WriteTimeout: 15 * time.Second,
		// IdleTimeout:  60 * time.Second,
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

	// Aguarda o sinal de interrupção
	sig := <-sigChan
	log.Printf("Sinal de desligamento recebido: %v, iniciando graceful shutdown...\n", sig)

	// Cria um contexto com timeout para o shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Ajuste o timeout conforme necessário
	defer cancel()

	// Inicia o processo de shutdown do servidor
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Erro durante o shutdown do servidor: %v\n", err)
	}

	log.Println("Servidor desligado com sucesso.")
}
