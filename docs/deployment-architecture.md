# Arquitetura de Deployment Escalável

## Diagrama de Arquitetura

```mermaid
graph TB
    subgraph "Load Balancer Layer"
        LB[Cloud Load Balancer]
    end

    subgraph "Kubernetes Cluster"
        subgraph "Ingress Layer"
            ING[Ingress Controller]
        end

        subgraph "Application Layer"
            subgraph "URL Shortener Service"
                APP1[URL Shortener Pod 1]
                APP2[URL Shortener Pod 2]
                APP3[URL Shortener Pod 3]
                HPA[Horizontal Pod Autoscaler]
            end
        end

        subgraph "Cache Layer"
            subgraph "Redis Cluster"
                REDIS1[Redis Master]
                REDIS2[Redis Replica 1]
                REDIS3[Redis Replica 2]
            end
        end

        subgraph "Database Layer"
            subgraph "PostgreSQL Cluster"
                PG_MASTER[PostgreSQL Master]
                PG_REPLICA1[PostgreSQL Replica 1]
                PG_REPLICA2[PostgreSQL Replica 2]
            end
        end

        subgraph "Monitoring Layer"
            PROM[Prometheus]
            GRAFANA[Grafana]
            K6[K6 Load Testing]
        end
    end

    %% Connections
    LB --> ING
    ING --> APP1
    ING --> APP2
    ING --> APP3
    HPA --> APP1
    HPA --> APP2
    HPA --> APP3

    APP1 --> REDIS1
    APP2 --> REDIS1
    APP3 --> REDIS1
    REDIS1 --> REDIS2
    REDIS1 --> REDIS3

    APP1 --> PG_MASTER
    APP2 --> PG_MASTER
    APP3 --> PG_MASTER
    PG_MASTER --> PG_REPLICA1
    PG_MASTER --> PG_REPLICA2

    PROM --> APP1
    PROM --> APP2
    PROM --> APP3
    PROM --> REDIS1
    PROM --> PG_MASTER
    GRAFANA --> PROM
    K6 --> ING

    %% Styling
    classDef primary fill:#4a90e2,stroke:#2171c7,color:white
    classDef secondary fill:#50e3c2,stroke:#2bb8a3,color:white
    classDef database fill:#f5a623,stroke:#d4880f,color:white
    classDef cache fill:#7ed321,stroke:#5cb315,color:white
    classDef monitoring fill:#9013fe,stroke:#6f0fc7,color:white

    class LB,ING primary
    class APP1,APP2,APP3,HPA secondary
    class PG_MASTER,PG_REPLICA1,PG_REPLICA2 database
    class REDIS1,REDIS2,REDIS3 cache
    class PROM,GRAFANA,K6 monitoring
```

## Componentes da Arquitetura

### 1. Load Balancer Layer
- **Cloud Load Balancer**: Distribui o tráfego entre os nós do cluster
- Suporte a SSL/TLS
- Health checks automáticos
- Auto-scaling baseado em demanda

### 2. Kubernetes Cluster
#### Ingress Layer
- **Ingress Controller**: Gerencia o roteamento de tráfego HTTP/HTTPS
- Configuração de regras de roteamento
- SSL termination
- Rate limiting

#### Application Layer
- **URL Shortener Service**: 
  - Múltiplos pods para alta disponibilidade
  - Horizontal Pod Autoscaler (HPA) para auto-scaling
  - Resource limits e requests configurados
  - Liveness e readiness probes

#### Cache Layer
- **Redis Cluster**:
  - Master-Replica setup para alta disponibilidade
  - Persistência de dados
  - Auto-failover
  - Cache invalidation automático

#### Database Layer
- **PostgreSQL Cluster**:
  - Master-Replica setup
  - Replicação síncrona
  - Auto-failover
  - Backup automático
  - Point-in-time recovery

#### Monitoring Layer
- **Prometheus**: Coleta métricas
- **Grafana**: Visualização e dashboards
- **K6**: Testes de carga automatizados

## Escalabilidade

### Horizontal Scaling
1. **Application Layer**:
   - Auto-scaling baseado em CPU/Memory
   - Pods distribuídos em múltiplos nós
   - Load balancing automático

2. **Cache Layer**:
   - Redis Cluster com sharding
   - Replicação para leitura
   - Cache distribuído

3. **Database Layer**:
   - Read replicas para distribuir carga de leitura
   - Connection pooling
   - Query optimization

### Vertical Scaling
- Aumento de recursos (CPU/Memory) por pod
- Otimização de configurações JVM/GC
- Ajuste de resource limits

## Alta Disponibilidade

1. **Multi-AZ Deployment**:
   - Distribuição em múltiplas zonas de disponibilidade
   - Auto-failover entre zonas
   - Data replication entre zonas

2. **Disaster Recovery**:
   - Backup automático
   - Point-in-time recovery
   - Failover automático

## Monitoramento e Observabilidade

1. **Métricas**:
   - Latência
   - Throughput
   - Error rates
   - Resource utilization

2. **Logging**:
   - Centralized logging
   - Log aggregation
   - Log analysis

3. **Alerting**:
   - Proactive alerts
   - On-call rotation
   - Escalation policies

## Segurança

1. **Network Security**:
   - Network policies
   - Service mesh
   - TLS encryption

2. **Access Control**:
   - RBAC
   - Service accounts
   - Secret management

3. **Compliance**:
   - Audit logging
   - Security scanning
   - Compliance monitoring
``` 