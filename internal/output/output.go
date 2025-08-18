package output

import (
	"fmt"
	"os"
)

func Handle(format string, content string, stdout bool, file string) {
	if stdout {
		fmt.Println(content)
	}
	if file != "" {
		if err := os.WriteFile(file, []byte(content), 0644); err != nil {
			fmt.Printf("Failed to write %s file: %v\n", format, err)
		}
	}
}
