<div align="center">

# DevForge

### Plataforma DevSecOps de ponta a ponta — da aplicação ao MLOps

Um encurtador de URLs em Go é o pretexto; o valor está na **infraestrutura DevSecOps** ao redor dele: provisionamento, orquestração, segurança automatizada, observabilidade e um serviço de machine learning operado pela mesma plataforma.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Python-3776AB?style=flat-square&logo=python&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white)
![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=flat-square&logo=kubernetes&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-7B42BC?style=flat-square&logo=terraform&logoColor=white)
![Vault](https://img.shields.io/badge/Vault-FFEC6E?style=flat-square&logo=vault&logoColor=black)
![Ansible](https://img.shields.io/badge/Ansible-EE0000?style=flat-square&logo=ansible&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat-square&logo=github-actions&logoColor=white)
![New Relic](https://img.shields.io/badge/New_Relic-1CE783?style=flat-square&logo=new-relic&logoColor=black)

</div>

---

## Visão geral

O DevForge demonstra o ciclo completo de engenharia de plataforma moderna. Uma API simples (encurtar e resolver URLs) é envolvida por uma stack de produção: o servidor é endurecido com Ansible, a infraestrutura é provisionada como código, os segredos vivem num cofre, a aplicação roda num cluster Kubernetes com autoscaling, um pipeline de CI/CD aplica três camadas de análise de segurança antes de cada entrega, a observabilidade é instrumentada em runtime, e um serviço de machine learning classifica o risco das URLs em tempo real.

O projeto é organizado em **8 fases**, cada uma resolvendo um problema real de operação.

---

## Arquitetura

```
                              ┌─────────────────────────────────────────┐
                              │            Cluster Kubernetes (Kind)     │
                              │                                          │
   Cliente ──HTTP──▶ Ingress ─┼──▶ Service(api) ──▶ Pod API (Go) ×2      │
                              │                      │     │             │
                              │                      │     ▼             │
                              │                      │   Service(redis)  │
                              │                      │     └▶ Pod Redis  │
                              │                      ▼                   │
                              │            Service(ml-service)           │
                              │              └▶ Pod ML (FastAPI) ×2      │
                              │                  (classificador phishing)│
                              │                                          │
                              │   HPA escala API e ML por CPU (2→5)       │
                              └──────────┬───────────────────────────────┘
                                         │
                          busca chave AES│ (runtime)
                                         ▼
                                  ┌─────────────┐
                                  │    Vault    │  cofre de segredos
                                  └─────────────┘

   Observabilidade: agent New Relic embutido na API (APM)
   Entrega: GitHub Actions → gosec → Trivy → Docker Hub → OWASP ZAP
```

### Fluxo de uma requisição

```
POST /shorten {"url": "http://paypa1-login.verify.xyz/account"}
   │
   ├─ API gera código curto, criptografa a URL (AES-256-GCM, chave do Vault)
   ├─ salva no Redis
   ├─ consulta o ml-service: a URL é suspeita?
   │
   ▼
{"code": "8232b348", "risk": "high", "score": 0.94}
```

---

## As 8 fases

| # | Fase | O que demonstra |
|---|------|-----------------|
| 1 | **API em Go** | Arquitetura em camadas, criptografia AES-256-GCM, Redis |
| 2 | **Hardening (Ansible)** | Idempotência, roles, endurecimento de SSH/firewall |
| 3 | **IaC (Terraform)** | Infraestrutura como código, state, registry de imagens |
| 4 | **Vault** | Gestão de segredos, seal/unseal (Shamir), produção |
| 5 | **Kubernetes (Kind)** | Deployment, Service, Ingress, HPA, Secrets |
| 6 | **CI/CD + DevSecOps** | gosec (SAST), Trivy (scan), OWASP ZAP (DAST) |
| 7 | **Observabilidade** | New Relic APM, transactions, distributed tracing |
| 8 | **MLOps** | Classificador de phishing (treino, serving, deploy, integração) |

---

## Stack

**Aplicação:** Go (API), Python/FastAPI (serviço ML), Redis
**Infraestrutura:** Docker, Kubernetes (Kind), Terraform, Ansible
**Segurança:** HashiCorp Vault, gosec, Trivy, OWASP ZAP, AES-256-GCM
**CI/CD:** GitHub Actions, Docker Hub
**Observabilidade:** New Relic
**ML:** scikit-learn (Random Forest)

---

## Estrutura do repositório

```
devforge/
├── app/                      # API em Go
│   ├── cmd/api/              # entrypoint
│   ├── internal/             # crypto, handler, service, repository
│   └── pkg/                  # clients (vault, mlclient)
├── ml-service/               # serviço de ML em Python
│   ├── train.py              # treino do modelo
│   ├── app/                  # FastAPI serving
│   └── model/                # modelo treinado (.pkl)
├── infra/
│   ├── ansible/              # hardening
│   ├── terraform/            # infra como código
│   └── kind/                 # config do cluster
├── k8s/base/                 # manifests (api, redis, ml-service, ingress)
└── .github/workflows/        # pipeline CI/CD
```

---

## Como rodar

### Pré-requisitos
Docker, kind, kubectl. (Go e Python rodam dentro de containers, não precisam estar instalados.)

### Pipeline de CI/CD
O pipeline roda automaticamente a cada push: testes, análise estática (gosec), scan de imagem (Trivy), push pro Docker Hub e DAST (OWASP ZAP). Veja `.github/workflows/ci.yml`.

### Cluster local
```bash
# cria o cluster
kind create cluster --config infra/kind/cluster.yaml

# sobe o Vault (ver runbook em infra/vault/README.md)
# e o Redis

# deploya a aplicação e o serviço de ML
kubectl create namespace devforge
kubectl apply -f k8s/base/api/
kubectl apply -f k8s/base/redis/
kubectl apply -f k8s/base/ml-service/
kubectl apply -f k8s/base/ingress/
```

### Treinar o modelo de ML
```bash
cd ml-service
docker run --rm -v "$PWD":/app -w /app python:3.12-slim \
  sh -c "pip install -q scikit-learn numpy && python train.py"
```

---

## Decisões técnicas

**Por que Kind e não Kubernetes gerenciado?** Custo zero e demonstração de domínio real de K8s — manifests escritos à mão, não cliques num painel. Quem opera Kind entende o Kubernetes por dentro.

**Por que o segredo no Vault e não em variável de ambiente?** A chave de criptografia começou hardcoded no código, foi pra um cofre com unseal por Shamir. A aplicação busca a chave em runtime, autenticada — o segredo nunca toca o repositório.

**Por que três scanners de segurança no pipeline?** Cada um cobre uma camada: gosec analisa o código (SAST), Trivy a imagem e suas dependências (scan de CVE, com gate que **bloqueia** o deploy), OWASP ZAP a aplicação rodando (DAST). Segurança que não bloqueia nada é decorativa.

**Por que MLOps e não "chamar uma API de IA"?** O classificador de phishing reusa toda a infraestrutura existente (Deployment, Service, HPA, CI/CD) — é a plataforma operando uma carga de ML como opera qualquer serviço. A IA degrada com elegância: se o serviço cair, o encurtador continua funcionando.

---

## Status e próximos passos

O núcleo está completo e funcional. Melhorias mapeadas: testes de cobertura mais ampla (hoje há testes no módulo de criptografia), CD automático (o deploy no cluster é manual), Vault dentro do cluster com Kubernetes Auth, dataset real de phishing e model registry para o serviço de ML, e regras de alerta no New Relic.

---

<div align="center">

Construído como projeto de portfólio de engenharia DevOps/Cloud.

</div>
