package main

import (
	"database/sql"
	"fmt"
	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/DiegoOpenheimer/go/clean-arch/configs"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/event/handler"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/graph"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/grpc/pb"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/grpc/service"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/web/webserver"
	"github.com/DiegoOpenheimer/go/clean-arch/pkg/events"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"os"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(config.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err)
	}
	executeMigration(db, config.DBDriver)
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	_ = eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	server := webserver.NewWebServer(config.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	server.AddHandler("/order", webOrderHandler.Create)
	fmt.Println("Starting web server on port", config.WebServerPort)
	go server.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		_ = grpcServer.Serve(lis)
	}()

	srv := graphqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	_ = http.ListenAndServe(":"+config.GraphQLServerPort, nil)
}

func executeMigration(db *sql.DB, driverName string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/internal/infra/database/migrations", path),
		driverName,
		driver,
	)
	if err != nil {
		panic(err)
	}
	_ = m.Steps(1)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
