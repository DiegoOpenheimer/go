package ports

import "context"

type UseCases[Input any, Output any] interface {
	Execute(ctx context.Context, input Input) Output
}
