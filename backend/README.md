# Backend Service

Este serviço de backend é uma aplicação Go responsável por gerar, armazenar e fornecer dados de séries temporais.

## Arquitetura

O backend utiliza as seguintes tecnologias:
- **Go**: Linguagem de programação principal.
- **net/http**: Para criar o servidor web e manipular as rotas da API.
- **Ent**: Um framework ORM para Go, usado para interagir com o banco de dados.
- **SQLite**: O banco de dados usado para armazenar os dados.

O serviço consiste em:
- Um servidor HTTP que expõe endpoints da API REST.
- Uma goroutine em segundo plano (`StartValueUpdater`) que gera e salva novos pares de valores no banco de dados a cada segundo.
- Um banco de dados SQLite (`radare.db`) que armazena o histórico de valores.

## Como Executar

Para executar o backend localmente para desenvolvimento:
```bash
cd backend
go run main.go
```
O servidor será iniciado na porta 8080.

Para executar usando Docker:
```bash
docker-compose up backend
```

## API Endpoints

### 1. `GET /current-values`

Retorna o registro de valor mais recente do banco de dados.

- **Resposta de Sucesso (200 OK):**
  ```json
  {
    "id": 123,
    "value1": 100,
    "value2": 50,
    "created_at": "2023-10-27T10:00:00Z"
  }
  ```

- **Resposta de Erro (404 Not Found):**
  Retornada se o banco de dados estiver vazio.
  ```
  Nenhum valor encontrado ainda. Tente novamente em alguns segundos.
  ```

### 2. `GET /values/history`

Retorna os últimos 10 registros de valor do banco de dados.

- **Resposta de Sucesso (200 OK):**
  ```json
  [
    {
      "id": 123,
      "value1": 100,
      "value2": 50,
      "created_at": "2023-10-27T10:00:00Z"
    },
    {
      "id": 122,
      "value1": 50,
      "value2": 100,
      "created_at": "2023-10-27T09:59:59Z"
    }
  ]
  ```

### 3. `GET /healthz`

Um endpoint de verificação de saúde.

- **Resposta de Sucesso (200 OK):**
  ```json
  {
    "status": "ok"
  }
  ```

## Banco de Dados

O serviço utiliza um banco de dados SQLite, e o arquivo do banco de dados (`radare.db`) é criado automaticamente no diretório raiz do backend.

### Schema (Tabela `valuelogs`)

O esquema é gerenciado pelo Ent e definido em `ent/schema/valuelog.go`.

| Coluna | Tipo | Descrição |
| --- | --- | --- |
| `id` | INTEGER | Chave primária, autoincremento. |
| `value1`| INTEGER | O primeiro valor inteiro. |
| `value2`| INTEGER | O segundo valor inteiro. |
| `created_at` | DATETIME | O timestamp de quando o registro foi criado. |