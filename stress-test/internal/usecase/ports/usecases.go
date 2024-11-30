package ports

type UseCases[Input any, Output any] interface {
	Execute(input Input) (Output, error)
}
