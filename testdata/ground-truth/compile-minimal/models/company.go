package models

type Company struct {
	ID        uint `morphe:"mandatory"`
	Name      string
	TaxID     string
	PersonIDs []uint
	Persons   []Person
}

func (m Company) GetIDName() CompanyIDName {
	return CompanyIDName{
		Name: m.Name,
	}
}

func (m Company) GetIDPrimary() CompanyIDPrimary {
	return CompanyIDPrimary{
		ID: m.ID,
	}
}
