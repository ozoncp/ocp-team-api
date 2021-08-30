package models

import "fmt"

type Team struct {
	Id          uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IsDeleted   bool   `db:"is_deleted"`
}

func (t Team) String() string {
	return fmt.Sprintf("{Id: %d, Name: %s, Description: %s, IsDeleted: %t}", t.Id, t.Name, t.Description, t.IsDeleted)
}
