package cli

import (
	"fmt"
	"github.com/DiegoOpenheimer/go/stress-test/configs"
	"github.com/DiegoOpenheimer/go/stress-test/internal/infra"
	"github.com/DiegoOpenheimer/go/stress-test/internal/usecase"
	"github.com/DiegoOpenheimer/go/stress-test/pgk/utils"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type RunnableCommand func(cmd *cobra.Command, args []string)

var rootCommand *cobra.Command

func NewRootCmd(useCase usecase.StressTestUseCase) *cobra.Command {
	return &cobra.Command{
		Use:   "start-stress",
		Short: "Load testing on web services",
		Long: `stress_test is a command-line tool built in Go, designed to efficiently perform load testing on web services.

The process is straightforward: the user provides the service URL, the total number of requests to be made, and the number of concurrent calls to execute. Based on these inputs, the system performs the load test and generates a detailed report at the end, including:
	•	Total time taken to complete the test.
	•	Total number of requests executed.
	•	Count of successful requests with HTTP 200 status.
	•	Detailed distribution of other HTTP status codes.

This tool is perfect for evaluating the performance and resilience of web services under different load scenarios.`,
		Run: runCommand(useCase),
	}
}

func runCommand(useCase usecase.StressTestUseCase) RunnableCommand {
	cfg := configs.GetConfig()
	return func(cmd *cobra.Command, args []string) {
		bar := progressbar.NewOptions(cfg.Requests,
			progressbar.OptionSetWriter(ansi.NewAnsiStdout()), //you should install "github.com/k0kubun/go-ansi"
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetDescription(fmt.Sprintf("Running %s", cfg.Url)),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}))
		response, err := useCase.Execute(usecase.StressTestInput{
			URL:         cfg.Url,
			Requests:    cfg.Requests,
			Concurrency: cfg.Concurrency,
			OnProgress: func(progress int) {
				_ = bar.Add(progress)
			},
		})
		utils.CheckError(err)
		fmt.Println("\n\nResult:")
		table := tablewriter.NewWriter(os.Stdout)
		var headers []string
		var values []string
		headers = append(headers, "totalTime", "totalRequests")
		values = append(values, fmt.Sprintf("%v", response["totalTime"]), fmt.Sprintf("%v", response["totalRequests"]))
		for k, v := range response {
			if k != "totalTime" && k != "totalRequests" {
				headers = append(headers, k)
				values = append(values, fmt.Sprintf("%v", v))
			}
		}
		table.SetHeader(headers)
		table.Append(values)
		table.Render()
	}
}

func Execute() {
	err := rootCommand.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(configs.LoadConfig)

	rootCommand = NewRootCmd(usecase.NewStressTest(infra.NewHttpWebService()))

	rootCommand.PersistentFlags().StringP("url", "u", "", "URL (required)")
	rootCommand.PersistentFlags().IntP("requests", "r", 100, " Total requests")
	rootCommand.PersistentFlags().IntP("concurrency", "c", 10, "Number of concurrent requests")
	_ = viper.BindPFlag("url", rootCommand.PersistentFlags().Lookup("url"))
	_ = viper.BindPFlag("requests", rootCommand.PersistentFlags().Lookup("requests"))
	_ = viper.BindPFlag("concurrency", rootCommand.PersistentFlags().Lookup("concurrency"))
}
