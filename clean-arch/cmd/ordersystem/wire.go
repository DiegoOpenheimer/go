//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/DiegoOpenheimer/go/clean-arch/internal/entity"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/event"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/database"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/infra/web"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/usecase"
	"github.com/DiegoOpenheimer/go/clean-arch/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(
	db *sql.DB,
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrderUseCase *usecase.ListOrdersUseCase,
) *web.OrderHandler {
	wire.Build(web.NewWebOrderHandler)
	return &web.OrderHandler{}
}
