# Dockerfile

# Primeira etapa - Compilação
FROM golang:alpine AS builder

# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar o código fonte para o contêiner
COPY . .

# Baixar as dependências e compilar o programa
RUN apk update && \
    apk add --no-cache git && \
    go mod download && \
    go build ./cmd/ordersystem

# Segunda etapa - Imagem final
FROM alpine:latest

# Copiar o executável do diretório de compilação para a imagem final
COPY --from=builder /app/ordersystem /usr/local/bin/ordersystem
COPY --from=builder /app/internal/infra/database/migrations /internal/infra/database/migrations

# Definir argumentos padrão
ENTRYPOINT ["ordersystem"]
