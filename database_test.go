package golang_database

import (
	"testing"
)

func TestEmpty(t *testing.T) {

}

func TestOpenConnection(t *testing.T) {
	db := GetConnection()
	defer db.Close()
}
