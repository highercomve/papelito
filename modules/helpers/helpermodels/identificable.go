package helpermodels

import (
	"github.com/highercomve/papelito/modules/helpers/prnx"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Identificable interface of something with ID
type Identificable interface {
	GetID() string
	GetPrn() string
	SetPrn(service string) Identificable
}

// Identification identification fields for a model
type Identification struct {
	ID  string `json:"id" bson:"_id"`
	PRN string `json:"prn" bson:"prn"`
}

// GetID get the id of a Identification
func (identification *Identification) GetID() string {
	return identification.ID
}

// GetPrn get the id of a Identification
func (identification *Identification) GetPrn() string {
	return identification.PRN
}

// GetID get the id of a Identification
func (identification *Identification) SetPrn(service string) Identificable {
	identification.PRN = prnx.IDGetPrn(identification.ID, service)

	return identification
}

// Ownable interface of something with ID
type Ownable interface {
	GetOwnerID() string
	GetOwnerPrn() string
	SetOwnerPrn(service string) Ownable
}

// Ownership define owner of a resource
type Ownership struct {
	OwnerID  string `json:"owner_id" bson:"owner_id"`
	OwnerPrn string `json:"owner" bson:"owner"`
}

// GetOwnerID get the id of a owner
func (owner *Ownership) GetOwnerID() string {
	return owner.OwnerID
}

// GetOwnerPrn get the prn of a owner
func (owner *Ownership) GetOwnerPrn() string {
	return owner.OwnerPrn
}

// SetOwnerPrn set the prn of a owner
func (owner *Ownership) SetOwnerPrn(service string) Ownable {
	owner.OwnerPrn = prnx.IDGetPrn(owner.OwnerID, "accounts")
	return owner
}

// NewIdentification Create new indentification
func NewIdentification(service string) Identification {
	identification := Identification{
		ID: primitive.NewObjectID().Hex(),
	}
	identification.PRN = prnx.IDGetPrn(identification.ID, service)

	return identification
}
