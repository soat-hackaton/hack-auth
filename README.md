# 🔐 Hack Auth Service

Microsserviço de Autenticação e Gestão de Usuários para o Hackathon SOAT.
Desenvolvido em **Go** com arquitetura limpa, utilizando **DynamoDB** e deploy via **Kubernetes (EKS)**.

![Go Version](https://img.shields.io/badge/go-1.24-blue)
![Coverage](https://img.shields.io/badge/coverage-95%25-brightgreen)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)

## 🚀 Tecnologias

- **Linguagem:** Go 1.24
- **Framework Web:** Gin Gonic
- **Banco de Dados:** Amazon DynamoDB
- **Autenticação:** JWT (JSON Web Token) + BCrypt
- **Infraestrutura:** Docker, Kubernetes (EKS), Terraform
- **CI/CD:** GitHub Actions

## ⚙️ Configuração Local

### Pré-requisitos
- Go 1.24+
- Docker ` Docker Compose
- Acesso à AWS (para DynamoDB)

### Variáveis de Ambiente
Copie o exemplo e configure suas credenciais:
```bash
cp .env-example .env
```
*Nota: Se estiver usando AWS Academy, lembre-se de atualizar `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` e `AWS_SESSION_TOKEN` a cada 4 horas.*

### Rodando com Docker (Recomendado)
O ambiente de desenvolvimento utiliza **Hot Reload** (Air). Qualquer alteração no código reinicia o servidor automaticamente.

```bash
docker compose up --build
```
A API estará disponível em: `http://localhost:8080`

## 🧪 Testes e Qualidade

O projeto possui uma barreira de qualidade no CI/CD que exige no mínimo **80% de cobertura** de testes.

Para rodar os testes localmente e ver a cobertura:
```bash
# Roda testes unitários nos pacotes de domínio
go test -v -coverprofile=coverage.out ./internal/service/... ./internal/handler/...

# Visualiza o relatório no terminal
go tool cover -func=coverage.out

# (Opcional) Visualiza em HTML no navegador
go tool cover -html=coverage.out
```

## 🔌 API Endpoints

### 1. Criar Usuário (SignUp)
**POST** `/auth/signup`

```json
{
  "name": "João da Silva",
  "email": "joao@email.com",
  "password": "senha_super_secreta"
}
```

### 2. Login
**POST** `/auth/login`

```json
{
  "email": "joao@email.com",
  "password": "senha_super_secreta"
}
```
**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## 📦 Deploy e Infraestrutura

O deploy é automatizado via **GitHub Actions** para o cluster EKS na AWS.

### 🔄 Pipeline de CI/CD (DevSecOps)

A nossa pipeline reflete as melhores práticas, com validações rigorosas e deploy automatizado:

1. **Testes e Qualidade (CI)**
   - Execução dos testes unitários com barreira obrigatória (mínimo de **80% de cobertura**).
   - Validação de segurança e qualidade do código fonte de forma automatizada através do **SonarQube/SonarCloud**.
2. **Automação de Pull Request**
   - Após rodar os testes com sucesso na branch `dev`, um PR é aberto de forma automática para a branch `main`.
3. **Build e Deploy (CD)**
   - Executado **somente** na branch `main` (ou manualmente em casos de bypass).
   - Criação automatizada da imagem Docker com envio para o **Amazon ECR**.
   - Injeção segura de credenciais e atualização dos manifestos usando Kustomize.
   - Deploy automático no cluster **AWS EKS**.

### 🛡️ Boas Práticas Adotadas

- **Proteção de Branch (`main`):** Commits diretos não são permitidos; todo fluxo exige pull requests que devem respeitar as regras do CI (cobertura + Sonar).
- **Gestão de Secrets:** As credenciais sensíveis nunca são versionadas nem salvas no repo. Elas ficam no GitHub Secrets e são injetadas em tempo de runtime no Kubernetes via env vars/secrets.

### Estrutura K8s
O projeto utiliza **Kustomize** para gestão de ambientes:
- `k8s/base`: Definições comuns (Deployment, Service, HPA).
- `k8s/overlays/dev`: Customizações de ambiente (Réplicas, Tags de Imagem).