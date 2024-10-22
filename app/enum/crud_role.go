package enum

type CrudRole string

const (
	CRUD_ROLE_CREATE CrudRole = "create"
	CRUD_ROLE_READ   CrudRole = "read"
	CRUD_ROLE_UPDATE CrudRole = "update"
	CRUD_ROLE_DELETE CrudRole = "delete"
)

func GetCrudRole(t CrudRole) CrudRole {
	switch t {
	case CRUD_ROLE_CREATE:
		return CRUD_ROLE_CREATE
	case CRUD_ROLE_READ:
		return CRUD_ROLE_READ
	case CRUD_ROLE_UPDATE:
		return CRUD_ROLE_UPDATE
	default:
		return CRUD_ROLE_DELETE
	}
}
