# URL Shortener Service

Serviço de encurtamento de URLs com expiração automática, desenvolvido em Go.

## Características

- Encurtamento de URLs com códigos únicos
- Expiração automática de URLs após 24 horas
- Redirecionamento automático para URLs originais
- Validação de URLs e protocolos
- Métricas e monitoramento com Prometheus e Grafana
- Armazenamento em Redis
- API RESTful
- Testes unitários
- Documentação completa

## Arquitetura

O projeto segue a arquitetura CMD (Command Query Responsibility Segregation) com as seguintes camadas:

```
.
├── cmd/                    # Pontos de entrada da aplicação
│   └── server/            # Servidor HTTP
├── internal/              # Código privado da aplicação
│   ├── api/              # Handlers HTTP e rotas
│   ├── domain/           # Entidades e regras de negócio
│   ├── repository/       # Camada de persistência
│   ├── service/          # Lógica de negócio
│   └── metrics/          # Métricas e monitoramento
├── pkg/                   # Código público reutilizável
└── test/                 # Testes de integração
```

## Endpoints da API

### 1. Encurtar URL
```bash
POST /shorten
Content-Type: application/json

{
    "url": "https://www.example.com"
}
```

Resposta:
```json
{
    "short_url": "http://url.li/Ab3Cd4Ef"
}
```

### 2. Redirecionar para URL Original
```bash
GET /:shortURL
```

Redireciona automaticamente para a URL original.

### 3. Obter Informações da URL
```bash
GET /info/:shortURL
```

Resposta:
```json
{
    "short_url": "http://url.li/Ab3Cd4Ef",
    "original_url": "https://www.example.com",
    "expires_at": "2024-02-21T15:04:05Z"
}
```

### 4. Métricas
```bash
GET /metrics
```
Endpoint Prometheus com métricas do serviço.

### 5. Health Check
```bash
GET /health
```
Endpoint para verificação de saúde do serviço.

## Métricas Disponíveis

### Métricas HTTP
- `http_requests_total`: Total de requisições HTTP por método, endpoint e status
- `http_request_duration_seconds`: Duração das requisições HTTP em segundos

### Métricas do Serviço
- `url_shortening_total`: Total de URLs encurtadas
- `url_redirects_total`: Total de redirecionamentos
- `active_urls`: Número atual de URLs ativas

## Monitoramento

O serviço inclui integração com Prometheus e Grafana para monitoramento:

### Prometheus
- Endpoint: `http://localhost:9090`
- Configuração em `prometheus.yml`
- Coleta métricas a cada 15 segundos

### Grafana
- Interface: `http://localhost:3000`
- Login padrão: admin/admin
- Dashboards pré-configurados para:
  - Taxa de requisições
  - Latência média
  - URLs ativas
  - Total de URLs encurtadas
  - Taxa de sucesso
  - Redirecionamentos

## Requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- Redis

## Instalação

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/url-shortener.git
cd url-shortener
```

2. Instale as dependências:
```bash
go mod download
```

3. Inicie os serviços com Docker Compose:
```bash
docker-compose up -d
```

4. Execute a aplicação:
```bash
go run cmd/server/main.go
```

## Testes

Execute os testes unitários:
```bash
go test ./...
```

Execute os testes de integração:
```bash
go test ./test/...
```

### Testes de Carga com K6

O projeto inclui testes de carga usando K6. Para executar os testes:

1. Certifique-se de que o serviço está rodando:
```bash
docker-compose up -d
```

2. Execute o teste de carga:
```bash
docker-compose run k6 run /scripts/load-test.js
```

O teste de carga inclui:
- Simulação de carga progressiva (50-100 usuários virtuais)
- Teste de todos os endpoints da API
- Métricas de performance e erros
- Integração com Prometheus para visualização das métricas

#### Configurações do Teste
- Duração total: 9 minutos
- Estágios:
  - 1 min: Rampa de subida para 50 usuários
  - 3 min: Manter 50 usuários
  - 1 min: Aumentar para 100 usuários
  - 3 min: Manter 100 usuários
  - 1 min: Rampa de descida
- Thresholds:
  - 95% das requisições devem completar em menos de 500ms
  - Taxa de erro deve ser menor que 10%

#### Visualizando Resultados
Os resultados dos testes de carga são automaticamente enviados para o Prometheus e podem ser visualizados no Grafana:
1. Acesse o Grafana (http://localhost:3000)
2. Importe o dashboard de testes de carga
3. Visualize as métricas de performance

## Exemplo de Uso

1. Encurtar uma URL:
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com"}'
```

2. Obter informações da URL:
```bash
curl http://localhost:8080/info/Ab3Cd4Ef
```

3. Acessar a URL encurtada:
```bash
curl -L http://localhost:8080/Ab3Cd4Ef
```

## Configuração do Ambiente

O projeto utiliza variáveis de ambiente para configuração. Copie o arquivo `.env.example` para `.env` e ajuste as variáveis conforme necessário:

```bash
cp .env.example .env
```

### Variáveis de Ambiente

- `SERVER_PORT`: Porta do servidor HTTP (padrão: 8080)
- `REDIS_HOST`: Host do Redis (padrão: localhost)
- `REDIS_PORT`: Porta do Redis (padrão: 6379)
- `REDIS_PASSWORD`: Senha do Redis (opcional)
- `BASE_URL`: URL base para as URLs encurtadas (padrão: http://url.li)
- `URL_DURATION`: Duração de expiração das URLs (padrão: 24h)

## Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes. 