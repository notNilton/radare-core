# Frontend WebApp

Este é um aplicativo de página única (SPA) em React que visualiza dados de séries temporais de um serviço de backend.

## Arquitetura

O frontend é construído com as seguintes tecnologias:
- **React**: Biblioteca principal para a construção da interface do usuário.
- **Vite**: Ferramenta de build para desenvolvimento e produção rápidos.
- **TypeScript**: Para tipagem estática e um desenvolvimento mais robusto.
- **ReactFlow**: Para a interface de arrastar e soltar de nós e arestas.
- **Chart.js & PrimeReact**: Para exibir gráficos e componentes de UI.
- **Sass**: Para estilização.

### Estrutura dos Componentes

- **`main.tsx`**: O ponto de entrada da aplicação React.
- **`App.tsx`**: O componente raiz que monta o layout principal.
- **`components/`**: Contém todos os componentes React.
  - **`Canva/Node.tsx`**: O componente principal que renderiza a interface do ReactFlow. A funcionalidade original de "reconciliação" foi preservada, mas não está conectada ao backend atual.
  - **`Graph/GraphComponent.tsx`**: Exibe o histórico de valores em um gráfico de linhas, buscando dados do endpoint `/api/values/history`.
  - **`Sidebar/SidebarComponent.tsx`**: Exibe os valores mais recentes, buscando dados do endpoint `/api/current-values`.
- **`api/service.ts`**: Centraliza a lógica de comunicação com o backend.
- **`vite.config.ts`**: Contém a configuração do Vite, incluindo um proxy reverso para as chamadas de API para o backend.

## Como Executar

Para executar o frontend localmente para desenvolvimento:
```bash
cd webapp
pnpm install
pnpm dev
```
A aplicação estará disponível em `http://localhost:5173`.

Para executar usando Docker:
```bash
docker-compose up frontend
```
A aplicação será servida na porta 80.

## Fluxo de Dados

1. O `GraphComponent` e o `SidebarComponent` buscam dados do backend a cada 2 segundos usando as funções em `api/service.ts`.
2. As chamadas de API (para `/api/...`) são roteadas para o serviço de backend (`http://backend:8080`) pelo proxy do Vite.
3. O `GraphComponent` renderiza o histórico de dados usando `Chart.js`.
4. O `SidebarComponent` exibe os valores mais recentes.
5. O componente `Canva/Node.tsx` contém a lógica para uma interface de nós e arestas, mas sua funcionalidade de reconciliação não está totalmente integrada com o backend atual.

## Variáveis de Ambiente

- **`VITE_API_URL`**: A URL base para as chamadas de API. Se não for definida, o padrão é `/api`, que funciona com a configuração do proxy.