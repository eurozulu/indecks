package indecks

import "fmt"

type Property interface {
	fmt.Stringer
}

type PropertyParser interface {
	ParseProperties(p string, by []byte) (map[string]Property, error)
}
