package glob

import "strings"

// ContainsGlobChars checks if a path contains glob pattern characters
func ContainsGlobChars(path string) bool {
	return strings.ContainsAny(path, "*?[")
}
