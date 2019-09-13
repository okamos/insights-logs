package ezinsights

import (
	"fmt"

	"github.com/okamos/insights-logs/ui"
)

// Run command
func Run() int {
	err := ui.Draw(version)
	if err != nil {
		fmt.Print(err)
		return 1
	}
	return 0
}
