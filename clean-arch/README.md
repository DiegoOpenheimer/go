# Order System

Atividade de Clean Arch que utiliza Go, RabbitMQ, MySQL e GraphQL. Abaixo estão as instruções para configurar e executar o projeto utilizando Docker Compose.

## Pré-requisitos

- Docker
- Docker Compose

## Executando o Projeto

Para iniciar todos os serviços (MySQL, RabbitMQ e o sistema de pedidos), execute o comando abaixo:

```sh
docker-compose up --build

Pode usar o arquivo Makefile para executar o comando acima, basta executar o comando abaixo:
make run
```

Este comando irá:

- Construir a imagem Docker do sistema de pedidos.
- Iniciar os contêineres do MySQL e RabbitMQ.
- Iniciar o contêiner do sistema de pedidos, que depende dos serviços MySQL e RabbitMQ estarem saudáveis
  - Aguardar os logs do container ordersystem informando que os servidores web, gRPC e GraphQL estão ok.

## Acessando os Serviços

- **Web Server**: Acesse `http://localhost:8000`
- **gRPC Server**: Acesse `localhost:50051`
- **GraphQL Server**: Acesse `http://localhost:8080`
- **RabbitMQ Management**: Acesse `http://localhost:15672` (usuário: `guest`, senha: `guest`)

## Parando os Serviços

Para parar e remover todos os contêineres, execute:

```sh
docker-compose down
```

## Observações

- Certifique-se de que as portas `3306`, `5672`, `15672`, `8000`, `50051` e `8080` estejam livres no seu host.
- O arquivo `docker-compose.yaml` está configurado para mapear as portas dos contêineres para as portas do host.

## Executar pelo go

- Caso queira executar o projeto sem o docker, execute o comando abaixo:
- Necessário ter o go 1.23 instalado na máquina e o rabbitmq e mysql rodando (pode utilizar o docker compose para isso, apenas pare o serviço ordersystem).
```sh
go run cd ./cmd/ordersystem

```
Na raiz do projeto contém um arquivo de environment .env na qual pode ser editado conforme o uso.
   ```dotenv
   DB_DRIVER=mysql
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=orders
   WEB_SERVER_PORT=:8000
   GRPC_SERVER_PORT=50051
   GRAPHQL_SERVER_PORT=8080
   RABBITMQ_URL=amqp://guest:guest@localhost:5672/
   ```

## Migration

__Para fins didáticos, foi utilizado o migrate dentro do próprio projeto, então na execução deste irá rodar as migrations automaticamente.__

## Endpoints
Dentro da pasta `api` contém um arquivo order.http com os endpoints para teste.