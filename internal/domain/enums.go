package domain

type UserRoleEnum string

const (
	UserRoleEnumUnspecified UserRoleEnum = "UNSPECIFIED"
	UserRoleEnumTrader      UserRoleEnum = "TRADER"
	UserRoleEnumViewer      UserRoleEnum = "VIEWER"
	UserRoleEnumAdmin       UserRoleEnum = "ADMIN"
)

func (u UserRoleEnum) String() string {
	return string(u)
}

type UserRolesEnum []UserRoleEnum
