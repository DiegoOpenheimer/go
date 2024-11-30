package usecase

import (
	"errors"
	"fmt"
	"github.com/DiegoOpenheimer/go/stress-test/internal/usecase/ports"
	"github.com/go-playground/validator/v10"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type StressTestUseCase = ports.UseCases[StressTestInput, StressTestOutput]

type StressTestInput struct {
	URL         string `validate:"required,url"`
	Requests    int
	Concurrency int
	OnProgress  func(progress int)
}

type StressTestOutput map[string]any

type StressTest struct {
	webService ports.WebService
	mu         sync.Mutex
}

func NewStressTest(webService ports.WebService) *StressTest {
	return &StressTest{webService: webService}
}

func (st *StressTest) Execute(input StressTestInput) (StressTestOutput, error) {
	if err := st.validate(input); err != nil {
		return nil, err
	}
	stressTestOutput := make(StressTestOutput)
	start := time.Now()
	var wg sync.WaitGroup
	numberRequests := int64(input.Requests)
	for i := 0; i < input.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.worker(input, &stressTestOutput, &numberRequests)
		}()
	}
	wg.Wait()
	stressTestOutput["totalTime"] = time.Since(start).String()
	stressTestOutput["totalRequests"] = strconv.Itoa(input.Requests - int(numberRequests))
	return stressTestOutput, nil
}

func (st *StressTest) validate(input StressTestInput) error {
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				return errors.New(fmt.Sprintf("Error in the field '%s': It's required.", err.Field()))
			case "url":
				return errors.New(fmt.Sprintf("Error in the field '%s': It's should be a valid URL.", err.Field()))
			default:
				return errors.New(fmt.Sprintf("Error in the field '%s': %s.\n", err.Field(), err.Error()))
			}
		}
		return errors.New("validation error")
	}
	return err
}

func (st *StressTest) worker(input StressTestInput, stressTestOutput *StressTestOutput, numberRequests *int64) {
	for {
		if *numberRequests <= 0 {
			break
		}
		atomic.AddInt64(numberRequests, -1)
		result, err := st.webService.Request(input.URL)
		st.mu.Lock()
		if err != nil {
			if (*stressTestOutput)[err.Error()] == nil {
				(*stressTestOutput)[err.Error()] = 0
			}
			(*stressTestOutput)[err.Error()] = (*stressTestOutput)[err.Error()].(int) + 1
		} else {
			key := "Status " + strconv.Itoa(result.Status)
			if (*stressTestOutput)[key] == nil {
				(*stressTestOutput)[key] = 0
			}
			(*stressTestOutput)[key] = (*stressTestOutput)[key].(int) + 1
		}
		input.OnProgress(1)
		st.mu.Unlock()
	}
}
