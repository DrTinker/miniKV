package helper

import "fmt"

func RaftPrefix(id, term int, format string) string {
	return fmt.Sprintf("node: %d term: %d %s", id, term, format)
}
