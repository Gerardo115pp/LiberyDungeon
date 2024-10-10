package models

import "fmt"

func getGreatestCommonDivisor(a int, b int) int {
	for b != 0 {
		a, b = b, a%b
	}

	return a
}

type MediaSize struct {
	Width  int
	Height int
}

func (ms MediaSize) String() string {
	return fmt.Sprintf("%dx%d", ms.Width, ms.Height)
}

func (ms MediaSize) AspectRatio() string {
	greatest_common_divisor := getGreatestCommonDivisor(ms.Width, ms.Height)

	return fmt.Sprintf("%d:%d", ms.Width/greatest_common_divisor, ms.Height/greatest_common_divisor)
}
