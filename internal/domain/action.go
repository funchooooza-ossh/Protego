package domain

import "strings"

type Action string

const (
	Read    Action = "read"
	Write   Action = "write"
	Manage  Action = "manage"
	Unknown Action = "unknown"
)

var methodToAction = map[string]Action{
	"get":     Read,
	"head":    Read,
	"options": Read,
	"post":    Write,
	"put":     Write,
	"patch":   Write,
	"delete":  Write,
	// "admin": Manage, // future
}

// MapMethodToAction maps HTTP method to Action (read/write/...)
func MapMethodToAction(method string) Action {
	if action, ok := methodToAction[strings.ToLower(method)]; ok {
		return action
	}
	return Unknown
}
