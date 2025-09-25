# Radare Challenge - Aplicação de Visualização de Dados

Este projeto é uma aplicação web full-stack que consiste em um backend Go e um frontend React. O backend gera dados de séries temporais, e o frontend os visualiza em tempo real.

## Arquitetura

O projeto é dividido em dois serviços principais:

- **`backend/`**: Um serviço Go que expõe uma API REST para fornecer dados. Ele usa um banco de dados SQLite para persistir os dados.
- **`webapp/`**: Uma aplicação de página única (SPA) em React que consome a API do backend e exibe os dados em gráficos e painéis.

Ambos os serviços são orquestrados usando Docker e Docker Compose.

## Como Executar a Aplicação Completa

A maneira mais fácil de executar o projeto é com o Docker Compose, que construirá e executará os contêineres do backend e do frontend.

### Pré-requisitos

- Docker
- Docker Compose

### Passos

1. **Clone o repositório:**
   ```bash
   git clone <URL_DO_REPOSITORIO>
   cd radare-challenge
   ```

2. **Construa e inicie os serviços:**
   ```bash
   docker-compose up --build
   ```

3. **Acesse a aplicação:**
   - O frontend estará disponível em **`http://localhost:80`**.
   - O backend estará disponível na porta **`8080`**.

## Detalhes dos Serviços

Para mais informações sobre cada serviço, consulte os `README.md` específicos:

- **[Documentação do Backend](./backend/README.md)**
- **[Documentação do Frontend](./webapp/README.md)**