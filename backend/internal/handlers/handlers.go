// handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	// Import the middleware package (not needed here, but might be in other handlers)
)

// Estrutura para representar os valores
type CurrentValues struct {
	Value1 int `json:"value1"`
	Value2 int `json:"value2"`
}

var (
	currentValues CurrentValues
	mutex         sync.RWMutex // Mutex para garantir acesso seguro aos valores
)

// Inicializa a variação dos valores
func init() {
	go updateValues()
}

// updateValues atualiza os valores a cada segundo
func updateValues() {
	for {
		// Bloqueia o mutex para escrita
		mutex.Lock()
		// Alterna os valores entre 50 e 100
		if currentValues.Value1 == 50 {
			currentValues.Value1 = 100
			currentValues.Value2 = 50
		} else {
			currentValues.Value1 = 50
			currentValues.Value2 = 100
		}
		// Libera o mutex
		mutex.Unlock()

		// Aguarda 1 segundo antes de atualizar novamente
		time.Sleep(1 * time.Second)
	}
}

// GetCurrentValues retorna os valores atuais
func GetCurrentValues(w http.ResponseWriter, r *http.Request) error { // Retorna um erro
	// Configura o cabeçalho para indicar que a resposta é JSON
	w.Header().Set("Content-Type", "application/json")

	// Bloqueia o mutex para leitura
	mutex.RLock()
	values := currentValues
	mutex.RUnlock()

	// Converte a estrutura para JSON e envia como resposta
	if err := json.NewEncoder(w).Encode(values); err != nil {
		return err // Retorna o erro para o middleware
	}
	return nil
}

// HealthCheck retorna o status do servidor
func HealthCheck(w http.ResponseWriter, r *http.Request) error { // Retorna um erro
	// Você pode adicionar lógica aqui para verificar a saúde do seu servidor

	// Por enquanto, vamos apenas retornar um status 200 OK
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		return err
	}
	return nil
}
