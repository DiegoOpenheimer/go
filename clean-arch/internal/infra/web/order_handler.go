package web

import (
	"encoding/json"
	"net/http"

	"github.com/DiegoOpenheimer/go/clean-arch/internal/usecase"
)

type OrderHandler struct {
	createOrderUseCase *usecase.CreateOrderUseCase
	listOrderUseCase   *usecase.ListOrdersUseCase
}

func NewWebOrderHandler(
	createOrderUseCase *usecase.CreateOrderUseCase,
	listOrderUseCase *usecase.ListOrdersUseCase,
) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase,
		listOrderUseCase,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	output, err := h.createOrderUseCase.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) List(w http.ResponseWriter, _ *http.Request) {
	output, err := h.listOrderUseCase.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
