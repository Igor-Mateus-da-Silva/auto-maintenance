# AutoCare — Controle de Manutenção de Veículos

Sistema simples e direto para gestão de veículos, catálogo de manutenções preventivas/corretivas e acompanhamento de vencimentos por tempo, quilometragem ou ambos.

## Stack

- **Backend:** Go (biblioteca padrão `net/http` + `database/sql`), SQLite (`modernc.org/sqlite`), autenticação via JWT (`golang-jwt/jwt/v5`) em cookie `HttpOnly`.
- **Frontend:** Angular (Standalone Components, Signals, `@if`/`@for`, PWA), Tailwind CSS v4 (tema dark minimalista) e `@lucide/angular`.

---

## Como Executar Localmente

### 1. Backend (Go)

1. Entre na pasta `backend/`:
   ```bash
   cd backend
   ```
2. Copie o arquivo de exemplo de variáveis de ambiente:
   ```bash
   cp .env.example .env
   ```
   Variáveis principais (o backend carrega automaticamente o `.env` se presente e exige `APP_ACCESS_KEY` e `JWT_SECRET` definidos):
   - `PORT=8080`
   - `DB_PATH=./data/car_maintenance.db`
   - `CORS_ALLOWED_ORIGINS=http://localhost:4200`
   - `APP_ACCESS_KEY` *(chave utilizada para login no sistema)*
   - `JWT_SECRET` *(segredo aleatório com pelo menos 32 caracteres)*

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
2. Instale as dependências:
   ```bash
   npm install
   ```
3. Inicie o servidor de desenvolvimento:
   ```bash
   npm start
   ```
4. Acesse **`http://localhost:4200`** e faça login utilizando a chave definida em `APP_ACCESS_KEY`.

---

## Deploy

### Frontend na Vercel (PWA)

1. Importe o repositório na **Vercel** e defina o **Root Directory** como `frontend`.
2. Nas configurações do projeto (**Settings → Environment Variables**), adicione a variável:
   - `API_URL`: URL pública da sua API Go (ex.: `https://autocare-api.fly.dev/api`).
3. O arquivo `frontend/vercel.json` já configura o build (`npm run build`), o diretório de saída (`dist/frontend/browser`), rewrites de SPA e os headers de cache adequados para o Service Worker (`ngsw-worker.js` e `ngsw.json`).

### Backend em Produção (Docker)

O `backend/Dockerfile` utiliza build multi-stage (`CGO_ENABLED=0`) gerando uma imagem Alpine enxuta executando com usuário não-root e expondo a porta `8080`:

```bash
cd backend
docker build -t autocare-api .
docker run -d \
  -p 8080:8080 \
  -v autocare_data:/app/data \
  -e ENV=production \
  -e DB_PATH=/app/data/car_maintenance.db \
  -e APP_ACCESS_KEY="sua-chave-forte" \
  -e JWT_SECRET="seu-segredo-jwt-aleatorio-com-32-ou-mais-caracteres" \
  -e CORS_ALLOWED_ORIGINS="https://seu-app.vercel.app" \
  autocare-api
```

- **Persistência:** Monte sempre um volume persistente em `/app/data` (ex.: Fly Volumes, Railway Volume ou bind mount em VPS) para preservar o banco SQLite entre deploys.
- **CORS:** Informe em `CORS_ALLOWED_ORIGINS` a URL exata do frontend na Vercel (suporta múltiplas origens separadas por vírgula).
