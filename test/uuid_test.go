package test

import (
	"fmt"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestUUID(t *testing.T) {

	s := uuid.NewV4().String()
	fmt.Println(s)

}
