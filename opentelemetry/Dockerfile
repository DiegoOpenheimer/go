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
    go build -o service ./cmd/.

# Segunda etapa - Imagem final
FROM alpine:latest

# Copiar o executável do diretório de compilação para a imagem final
COPY --from=builder /app/service /usr/local/bin/service

# Definir argumentos padrão
ENTRYPOINT ["service"]
