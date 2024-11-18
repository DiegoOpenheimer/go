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
    go build ./cmd/deploy_cloud_run

# Segunda etapa - Imagem final
FROM alpine:latest

# Copiar o executável do diretório de compilação para a imagem final
COPY --from=builder /app/deploy_cloud_run /usr/local/bin/deploy_cloud_run

# Definir argumentos padrão
ENTRYPOINT ["deploy_cloud_run"]
