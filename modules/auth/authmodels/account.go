package authmodels

import (
	"github.com/highercomve/papelito/modules/helpers/helpermodels"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const AccountServicePrn = "accounts"

// AccountType Type of account
type AccountType string

// RolesType type of roles
type RolesType string

const (
	// AccountTypeAdmin define account type ADMIN
	AccountTypeAdmin = AccountType("ADMIN")

	// AccountTypeDevice define account type DEVICE
	AccountTypeDevice = AccountType("DEVICE")

	// AccountTypeOrg define account type ORG
	AccountTypeOrg = AccountType("ORG")

	// AccountTypeService define account type SERVICE
	AccountTypeService = AccountType("SERVICE")

	// AccountTypeUser define account type USER
	AccountTypeUser = AccountType("USER")

	// AccountTypeSessionUser define account type SESSION
	AccountTypeSessionUser = AccountType("SESSION")

	// AccountTypeClient define account type CLIENT
	AccountTypeClient = AccountType("CLIENT")

	// AccountTypeService define account type SERVICE
	AccountTypeResource = AccountType("RESOURCE")

	// RoleTypeUser role for users
	RoleTypeUser = RolesType("user")

	// RoleTypeService role for services
	RoleTypeService = RolesType("service")

	// RoleTypeResource role for resources
	RoleTypeResource = RolesType("resource")

	// RoleTypeAdmin role for admin actions
	RoleTypeAdmin = RolesType("admin")
)

// Account account information all the structure
type Account struct {
	helpermodels.Timestamp      `json:",inline" bson:",inline"`
	helpermodels.Identification `json:",inline" bson:",inline"`
	helpermodels.Ownership      `json:"-" bson:"-"`

	Email    string      `json:"email" bson:"email"`
	Nick     string      `json:"username" bson:"nick"`
	Password string      `json:"password" bson:"password"`
	Type     AccountType `json:"type" bson:"type"`
	Role     RolesType   `json:"role" bson:"role"`
}

func (acc *Account) GetServicePrn() string {
	return AccountServicePrn
}

// NewAccount create new account struct
func NewAccount() *Account {
	acc := &Account{
		Identification: helpermodels.Identification{
			ID: primitive.NewObjectID().Hex(),
		},
		Timestamp: helpermodels.Timestamp{},
	}

	acc.SetCreatedAt()
	acc.Type = AccountTypeUser
	acc.Role = RoleTypeUser
	return acc
}
