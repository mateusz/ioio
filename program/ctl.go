package program

import (
	"fmt"
	"strconv"
)

type ctl map[string]string

func (c ctl) string(key string) string {
	return c[key]
}

func (c ctl) int(key string) int {
	i, err := strconv.Atoi(c[key])
	if err != nil {
		fmt.Printf("Failed to convert ctl.%s: %s\n", key, err)
		return 0
	}
	return i
}

func (c ctl) float64(key string) float64 {
	f, err := strconv.ParseFloat(c[key], 64)
	if err != nil {
		fmt.Printf("Failed to convert ctl.%s: %s\n", key, err)
		return 0.0
	}
	return f
}
