package helpermodels

// Datable Db data
type Datable interface {
	Identificable
	Timeable
	Ownable

	GetServicePrn() string
}
