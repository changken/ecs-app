package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("Scheduled job starting...")

	apiKey := os.Getenv("API_KEY")

	report := map[string]any{
		"job":       "scheduled-report",
		"version":   os.Getenv("APP_VERSION"),
		"region":    os.Getenv("AWS_REGION"),
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		// 驗證 Secrets Manager 注入是否成功（不印出實際值）
		"secret_injected": apiKey != "",
	}

	b, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(b))
	log.Println("Job completed successfully, exiting 0")
}
