package idgenerator

import "github.com/teris-io/shortid"

// Generate is
func Generate() (id string) {
	id, _ = shortid.Generate()
	return
}
