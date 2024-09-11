package file_cleaner

import (
	"fmt"
	"time"
)

// refer: https://stackoverflow.com/a/45766707
func Timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("Time elapsed for %s: %v\n", name, time.Since(start))
	}
}
