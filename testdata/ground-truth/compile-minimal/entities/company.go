package entities

type Company struct {
	ID        uint
	Name      string
	TaxID     string
	PersonIDs []uint
	Persons   []Person
}

func (e Company) GetIDPrimary() CompanyIDPrimary {
	return CompanyIDPrimary{
		ID: e.ID,
	}
}
