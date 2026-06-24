# infra

Konfigurasi infrastruktur untuk deployment lokal, staging, dan production ERP UNSIA.

## Struktur

```
infra/
├── docker/
│   ├── docker-compose.local.yml    → Development (semua service + DB + broker)
│   ├── docker-compose.staging.yml
│   └── docker-compose.prod.yml
│
├── nginx/
│   ├── api-gateway.conf            → Routing ke semua service
│   ├── portal.conf                 → Frontend Next.js
│   └── services.conf
│
├── postgres/                       → Init config per DB modul
│   ├── core-db/
│   ├── reference-db/
│   ├── crm-db/
│   ├── pmb-db/
│   ├── finance-db/
│   ├── academic-db/
│   ├── hris-db/
│   ├── lms-db/
│   ├── assessment-db/
│   └── portal-db/
│
├── redis/
│   └── redis.conf
│
├── rabbitmq/
│   ├── definitions.json            → Exchange, queue, binding definitions
│   └── rabbitmq.conf
│
├── monitoring/
│   ├── prometheus/
│   ├── grafana/
│   ├── loki/                       → Log aggregation
│   └── alertmanager/
│
├── ci-cd/
│   ├── github-actions/
│   └── gitlab-ci/
│
├── backup/
│   ├── backup-postgres.sh
│   └── restore-postgres.sh
│
└── secrets/
    └── .env.example                → Template semua environment variable
```

## Quick Start (Local)

```bash
docker-compose -f docker/docker-compose.local.yml up -d
```

Service yang dijalankan: semua 10 PostgreSQL DB, Redis, RabbitMQ, semua Go services, Next.js portal.

## Port Default (Local)

| Service | Port |
|---------|------|
| API Gateway (Nginx) | 8080 |
| Portal Web | 3000 |
| Core Service | 8001 |
| Reference Service | 8002 |
| CRM Service | 8003 |
| PMB Service | 8004 |
| Finance Service | 8005 |
| Academic Service | 8006 |
| HRIS Service | 8007 |
| LMS Service | 8008 |
| Assessment Service | 8009 |
| Portal Service | 8010 |
| RabbitMQ Management | 15672 |
| Prometheus | 9090 |
| Grafana | 3001 |
