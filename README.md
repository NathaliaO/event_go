# Microserviço em Go, PostgreSQL e Docker

## 🚀 Como usar

- [Documentação Go](https://golang.org/doc/)
- [Documentação PostgreSQL](https://www.postgresql.org/docs/)
- [Docker Compose](https://docs.docker.com/compose/) 

### 1. Clonar e configurar
```bash
git clone https://github.com/NathaliaO/event_go.git
cd event_go

go mod tidy
```

### 2. Executar com Docker

#### 🍎 **Linux/Mac:**
```bash
make build
make run

```

#### 🪟 **Windows:**
```powershell
.\build.ps1 build
.\build.ps1 run

```

### 3. Acessar a aplicação - POSTMAN/INSOMNIA
- **API**: http://localhost:8080/health

## 📁 Estrutura do projeto

```
.
├── cmd/
│   └── server/
│       └── main.go              # Aplicação principal
├── internal/
│   ├── domain/                  # Estruturas de dados (User, Event)
│   ├── service/                 # Lógica de negócio
│   │   ├── *_service.go
│   │   └── *_service_test.go    # Testes unitários
│   ├── repository/              # Acesso a dados
│   │   ├── *_repository.go
│   │   └── *_repository_test.go # Testes unitários
│   └── handler/                 # Handlers HTTP
│       └── *_handler.go
├── Dockerfile                   # Configuração do container Go
├── docker-compose.yml           # Orquestração dos serviços
├── .dockerignore               # Arquivos ignorados pelo Docker
├── go.mod                      # Dependências Go
├── go.sum                      # Checksums das dependências
├── init.sql                    # Script de inicialização do banco
├── Makefile                    # Comandos úteis
└── README.md                   # Este arquivo
```

### Arquitetura

O projeto segue uma **arquitetura em camadas** (Clean Architecture):

1. **Handler** → Recebe requisições HTTP e valida dados
2. **Service** → Contém a lógica de negócio
3. **Repository** → Acessa o banco de dados
4. **Domain** → Estruturas de dados compartilhadas

**Nota**: A chave JWT é gerada automaticamente a cada inicialização da aplicação.

## 🌐 Endpoints da API

### Rotas públicas (sem autenticação)
- `GET /` - Página inicial
- `GET /health` - Status da API
- `POST /login` - Fazer login
- `POST /register` - Registrar novo usuário

### Rotas protegidas (requerem token JWT)
- `GET /users` - Listar usuários
- `POST /users` - Criar usuário
- `GET /profile` - Ver perfil do usuário logado

- `POST /api/events` - Recebe a lista de eventos
- `GET /api/stats/daily` - Retorna agregado por dia e site
## 🧪 Testes

### Executar testes

```bash
docker run --rm -v "$(pwd):/app" -w /app golang:1.21-alpine go test ./...
```
