// handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"radare.com/backend/ent"
	"radare.com/backend/ent/valuelog"
	"time"
)

// Handler é uma struct que contém o cliente do banco de dados.
type Handler struct {
	client *ent.Client
}

// New cria uma nova instância de Handler.
func New(client *ent.Client) *Handler {
	return &Handler{client: client}
}

// StartValueUpdater inicia uma goroutine para atualizar os valores no banco de dados.
func (h *Handler) StartValueUpdater() {
	var value1, value2 = 50, 100 // Valores iniciais
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Alterna os valores
			value1, value2 = value2, value1

			// Cria um novo registro no banco de dados
			_, err := h.client.ValueLog.
				Create().
				SetValue1(value1).
				SetValue2(value2).
				Save(context.Background())
			if err != nil {
				log.Printf("Erro ao salvar valores no banco de dados: %v", err)
			}
		}
	}
}

// GetCurrentValues retorna o registro mais recente do banco de dados.
func (h *Handler) GetCurrentValues(w http.ResponseWriter, r *http.Request) error {
	// Busca o último registro pelo ID em ordem decrescente
	latest, err := h.client.ValueLog.
		Query().
		Order(ent.Desc(valuelog.FieldID)).
		First(r.Context())
	if err != nil {
		// Se não houver registros, retorna um erro amigável
		if ent.IsNotFound(err) {
			http.Error(w, "Nenhum valor encontrado ainda. Tente novamente em alguns segundos.", http.StatusNotFound)
			return nil
		}
		// Para outros erros, retorna o erro para o middleware
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(latest); err != nil {
		return err
	}
	return nil
}

// GetValueHistory retorna os últimos 10 registros do banco de dados.
func (h *Handler) GetValueHistory(w http.ResponseWriter, r *http.Request) error {
	// Busca os últimos 10 registros
	history, err := h.client.ValueLog.
		Query().
		Order(ent.Desc(valuelog.FieldID)).
		Limit(10).
		All(r.Context())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		return err
	}
	return nil
}

// HealthCheck retorna o status do servidor.
func HealthCheck(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		return err
	}
	return nil
}