package entities

import (
	"github.com/kalo-build/dummy/enums"
)

type Person struct {
	Email       string
	ID          uint
	LastName    string
	Nationality enums.Nationality
	CompanyID   *uint
	Company     *Company
}

func (e Person) GetIDPrimary() PersonIDPrimary {
	return PersonIDPrimary{
		ID: e.ID,
	}
}
