package models

import (
	"fmt"
)

// Team is the representation of the team.
type Team struct {
	Id          uint64 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	IsDeleted   bool   `db:"is_deleted"`
}

// String is the method for converting Team struct to string representation.
func (t Team) String() string {
	return fmt.Sprintf("{Id: %d, Name: %s, Description: %s, IsDeleted: %t}", t.Id, t.Name, t.Description, t.IsDeleted)
}
