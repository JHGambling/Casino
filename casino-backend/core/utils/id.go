package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateID() string {
	// Just assume no error will happen :)
	id, _ := gonanoid.New()
	return id
}
