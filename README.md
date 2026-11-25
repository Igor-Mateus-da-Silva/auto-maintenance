# AutoCare — Controle de Manutenção de Veículos

Sistema simples e direto para gestão de veículos, catálogo de manutenções preventivas/corretivas e acompanhamento de vencimentos por tempo, quilometragem ou ambos.

## Stack

- **Backend:** Go (biblioteca padrão `net/http` + `database/sql`), SQLite (`modernc.org/sqlite`), autenticação via JWT (`golang-jwt/jwt/v5`).
- **Frontend:** Angular (Standalone Components, Signals, `@if`/`@for`), Tailwind CSS v4 (tema dark minimalista) e `lucide-angular`.

---

## Como Executar Localmente

### 1. Backend (Go)

1. Entre na pasta `backend/`:
   ```bash
   cd backend
   ```
2. Copie o arquivo de exemplo de variáveis de ambiente (ou defina no terminal):
   ```bash
   cp .env.example .env
   ```
   Variáveis principais (possuem valores padrão em desenvolvimento caso omitidas):
   - `PORT=8080`
   - `DB_PATH=./data/car_maintenance.db`
   - `CORS_ALLOWED_ORIGINS=http://localhost:4200`
   - `APP_ACCESS_KEY=admin123` *(chave utilizada para login no sistema)*
   - `JWT_SECRET=super-secret-jwt-key-auto-maintenance`

3. Inicie o servidor:
   ```bash
   go run ./cmd/api
   ```
   O servidor criará o banco SQLite automaticamente e aplicará as migrations na inicialização em `http://localhost:8080`.

4. Para rodar os testes do backend:
   ```bash
   go test -v ./...
   ```

---

### 2. Frontend (Angular)

1. Em outro terminal, entre na pasta `frontend/`:
   ```bash
   cd frontend
   ```
2. Instale as dependências (caso ainda não tenha instalado):
   ```bash
   npm install
   ```
3. Inicie o servidor de desenvolvimento:
   ```bash
   npm start
   ```
4. Acesse **`http://localhost:4200`** e faça login utilizando a chave definida em `APP_ACCESS_KEY` (padrão: `admin123`).
