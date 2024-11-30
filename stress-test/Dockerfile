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
    GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o stress-test ./cmd/.

# Segunda etapa - Imagem final
FROM alpine:latest

# Copiar o executável do diretório de compilação para a imagem final
COPY --from=builder /app/stress-test /usr/local/bin/stress-test

# Definir argumentos padrão
ENTRYPOINT ["stress-test"]
