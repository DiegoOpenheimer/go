package usecase

import (
	"github.com/DiegoOpenheimer/go/clean-arch/internal/entity"
	"github.com/DiegoOpenheimer/go/clean-arch/internal/usecase/dto"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() ([]dto.OrderOutputDTO, error) {
	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return nil, err
	}
	ordersOutput := make([]dto.OrderOutputDTO, 0)
	for _, order := range orders {
		ordersOutput = append(ordersOutput, dto.OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		})
	}
	return ordersOutput, nil
}
