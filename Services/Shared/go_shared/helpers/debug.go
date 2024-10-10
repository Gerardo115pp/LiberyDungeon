package helpers

import "fmt"

func PrintBytes(data []byte) {
	for _, b := range data {
		fmt.Printf("%#x ", b)
	}

	fmt.Println()
}
