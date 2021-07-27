package models

import "fmt"

type Team struct {
	Id uint64
	Name string
	Description string
}

func (t Team) String() string {
	return fmt.Sprintf("{Id: %d, Name: %s, Description: %s}", t.Id, t.Name, t.Description)
}