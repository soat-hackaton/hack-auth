# --- Estágio 1: Builder (Compilação) ---
FROM golang:1.25-alpine AS builder

# Instala certificados CA (necessário para o binário chamar HTTPS da AWS)
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copia e baixa dependências primeiro (Cache Layering)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário estático
# CGO_ENABLED=0 garante que não dependa de bibliotecas C do sistema
# -ldflags="-w -s" remove informações de debug para diminuir o tamanho
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o auth-api cmd/api/main.go

# --- Estágio 2: Runner (Imagem Final Minimalista) ---
FROM scratch

WORKDIR /app

# Copia os certificados SSL do estágio builder (CRUCIAL para DynamoDB/AWS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copia apenas o binário compilado
COPY --from=builder /app/auth-api .

# Expõe a porta
EXPOSE 8080

# Comando de entrada
CMD ["./auth-api"]