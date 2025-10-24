package utils

import "strconv"

// ParseInt64 parses the given string as a base-10 int64.
// It panics if the string cannot be parsed as an int64.
func Atoi64(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}
