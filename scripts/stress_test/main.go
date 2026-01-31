package main

import (
	"fmt"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/api/product",
	})
	attacker := vegeta.NewAttacker()

	var results vegeta.Results
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		results.Add(res)
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("95th percentile: %s\n", metrics.Latencies.P95)
	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Success ratio:   %.2f%%\n", metrics.Success*100)
	fmt.Printf("Status Codes:    %v\n", metrics.StatusCodes)

	// Also print a text report
	reporter := vegeta.NewTextReporter(&metrics)
	err := reporter.Report(os.Stdout)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
	}
}
