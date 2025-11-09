package domain

type CustomerType string

const (
	CustomerTypeIndividual CustomerType = "individual"
	CustomerTypeCompany    CustomerType = "company"
)

func (t CustomerType) String() string {
	return string(t)
}
