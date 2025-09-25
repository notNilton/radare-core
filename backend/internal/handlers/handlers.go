// handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"radare-core/backend/backend/ent"
	"radare-core/backend/backend/ent/valuelog"
	"time"
)

// Handler encapsula as dependências para os manipuladores HTTP, como o cliente do banco de dados.
type Handler struct {
	client *ent.Client
}

// New cria e retorna uma nova instância de Handler.
// Ele injeta o cliente do banco de dados para ser usado pelos manipuladores de rota.
func New(client *ent.Client) *Handler {
	return &Handler{client: client}
}

// StartValueUpdater é uma goroutine de longa duração que gera e salva novos valores no banco de dados
// a cada segundo. Ele alterna dois valores para simular dados de séries temporais.
func (h *Handler) StartValueUpdater() {
	// Define os valores iniciais que serão alternados.
	var value1, value2 = 50, 100
	// Cria um ticker que dispara a cada 1 segundo.
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// Aguarda o próximo tick.
		case <-ticker.C:
			// Alterna os valores.
			value1, value2 = value2, value1

			// Cria um novo registro de ValueLog no banco de dados com os novos valores.
			_, err := h.client.ValueLog.
				Create().
				SetValue1(value1).
				SetValue2(value2).
				Save(context.Background())
			if err != nil {
				// Loga qualquer erro que ocorra durante a inserção no banco de dados.
				log.Printf("Erro ao salvar valores no banco de dados: %v", err)
			}
		}
	}
}

// GetCurrentValues é um manipulador HTTP que busca e retorna o registro de valor mais recente.
func (h *Handler) GetCurrentValues(w http.ResponseWriter, r *http.Request) error {
	// Consulta o banco de dados para o último registro de ValueLog, ordenado por ID decrescente.
	latest, err := h.client.ValueLog.
		Query().
		Order(ent.Desc(valuelog.FieldID)).
		First(r.Context()) // First() retorna apenas o primeiro resultado.
	if err != nil {
		// Se nenhum registro for encontrado, retorna um status 404 Not Found.
		if ent.IsNotFound(err) {
			http.Error(w, "Nenhum valor encontrado ainda. Tente novamente em alguns segundos.", http.StatusNotFound)
			return nil // Retorna nil porque a resposta de erro já foi escrita.
		}
		// Para qualquer outro erro, retorna o erro para ser tratado pelo middleware.
		return err
	}

	// Define o cabeçalho Content-Type e codifica o resultado em JSON.
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(latest)
}

// GetValueHistory é um manipulador HTTP que busca e retorna os últimos 10 registros de valor.
func (h *Handler) GetValueHistory(w http.ResponseWriter, r *http.Request) error {
	// Consulta o banco de dados para os últimos 10 registros de ValueLog.
	history, err := h.client.ValueLog.
		Query().
		Order(ent.Desc(valuelog.FieldID)). // Ordena por ID para obter os mais recentes.
		Limit(10).                         // Limita o resultado a 10.
		All(r.Context())                   // All() retorna todos os resultados correspondentes.
	if err != nil {
		return err
	}

	// Define o cabeçalho Content-Type e codifica o resultado em JSON.
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(history)
}

// HealthCheck é um manipulador HTTP simples para verificações de saúde, retornando um status "ok".
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}