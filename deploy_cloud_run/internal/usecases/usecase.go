package usecases

type UseCase[Input any, Output any] interface {
	Execute(input Input) (Output, error)
}
