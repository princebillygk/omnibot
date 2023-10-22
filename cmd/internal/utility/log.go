package utility

import (
	"encoding/json"
	"log"
)

func DampVar(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	log.Print(string(b))
}
