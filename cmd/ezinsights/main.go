package main

import (
	"os"

	ezinsights "github.com/okamos/insights-logs"
)

func main() {
	os.Exit(ezinsights.Run())
}
