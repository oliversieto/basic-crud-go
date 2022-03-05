package database

import (
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	database, error := Connect()

	if error != nil {
		log.Fatal("Não foi possível conectar no banco de dados")
	}

	database.Close()
}
