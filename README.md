# Serviço de Rate Limiter

Um serviço de limitação de requisições (rate limiter) configurável, escrito em Go, que pode limitar requisições com base no endereço IP ou token de API.

## Funcionalidades

- Limitação por endereço IP
- Limitação por token de API (sobrepõe os limites de IP)
- Armazenamento em Redis com TTL configurável
- Limites e tempos de expiração configuráveis
- Fácil integração como middleware HTTP
- Suporte a Docker

## Requisitos

- Go 1.21 ou superior
- Redis 7.x
- Docker e Docker Compose (opcional)

## Início Rápido

### Usando Docker

1. Clone o repositório
2. Copie o arquivo de ambiente:
   ```bash
   cp .env.example .env
   ```
3. Inicie os serviços:
   ```bash
   docker-compose up -d
   ```

### Instalação Manual

1. Instale as dependências:
   ```bash
   go mod download
   ```
2. Configure as variáveis de ambiente (veja a seção Configuração)
3. Inicie o Redis
4. Execute a aplicação:
   ```bash
   go run main.go
   ```

## Configuração

O serviço pode ser configurado usando variáveis de ambiente ou um arquivo `.env`:

| Variável | Descrição | Valor Padrão |
|----------|-----------|--------------|
| REDIS_ADDR | Endereço do servidor Redis | localhost:6379 |
| REDIS_PASSWORD | Senha do Redis | "" |
| REDIS_DB | Número do banco de dados Redis | 0 |
| DEFAULT_LIMIT | Número padrão de requisições permitidas por janela | 10 |
| DEFAULT_EXPIRATION | Tamanho padrão da janela de tempo (segundos) | 60 |
| PORT | Porta do servidor | 8080 |
| LOG_LEVEL | Nível de log | info |

## Uso da API

### Limitação por IP

Por padrão, as requisições são limitadas por endereço IP. O IP do cliente é detectado automaticamente, com suporte a headers X-Forwarded-For.

Exemplo:
```bash
curl http://localhost:8080/hello
```

### Limitação por Token

Para usar a limitação por token, inclua o header API_KEY:

```bash
curl -H "API_KEY: seu-token" http://localhost:8080/hello
```

### Headers de Resposta

O serviço inclui informações de limite nas headers de resposta:
- X-RateLimit-Limit: Número máximo de requisições permitidas
- X-RateLimit-Remaining: Requisições restantes na janela atual
- X-RateLimit-Reset: Tempo em segundos até o limite ser resetado

### Resposta de Limite Excedido

Quando os limites são excedidos, o serviço retorna:
- Código de Status: 429
- Mensagem: "you have reached the maximum number of requests or actions allowed within a certain time frame"

## Testes

Execute os testes unitários:
```bash
go test ./...
```

Execute os testes de benchmark:
```bash
go test -bench=. -benchmem ./internal/middleware/
```

Execute os testes no Docker:
```bash
docker-compose up test
```

## Arquitetura

O serviço segue os princípios da arquitetura limpa:

- `internal/middleware`: Implementação do middleware HTTP
- `internal/limiter`: Lógica principal de limitação de requisições
- `internal/infra/database`: Implementações de armazenamento (Redis)
- `internal/logger`: Configuração de logs

A camada de armazenamento é abstraída através da interface `DatabaseStore`, facilitando a adição de novos backends de armazenamento.

## Como Contribuir

1. Faça um fork do repositório
2. Crie sua branch de feature
3. Faça commit das suas alterações
4. Faça push para a branch
5. Crie um Pull Request

## Licença

Este projeto está licenciado sob a Licença MIT.