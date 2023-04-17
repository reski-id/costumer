package utils

import (
	"errors"
	"strconv"
)

// StringToInt converts a string to an int
func StringToInt(str string) (int, error) {
	if str == "" {
		return 0, errors.New("empty string")
	}
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.New("invalid string")
	}
	return num, nil
}
