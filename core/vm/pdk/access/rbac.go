package access

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

var (
	RoleAdmin = ComputeRoleHash("ADMIN_ROLE") // RoleAdmin represents the admin role
)

type State struct {
	// role ==> account ==> hasRole
	Roles map[common.Hash]map[common.Address]bool
}

type RBAC struct {
	path *storage.Path
}

func NewRBAC(st *storage.Storage, fieldName string) *RBAC {
	return &RBAC{
		path: st.Field(fieldName),
	}
}

func (r *RBAC) HasRole(role common.Hash, account common.Address) bool {
	hasRolePath := r.path.Field("Roles").Map(role).Map(account)
	hasRole, err := storage.Get[bool](hasRolePath)
	if err != nil {
		return false
	}
	return hasRole
}

func (r *RBAC) GrantRole(role common.Hash, account common.Address) {
	rolePath := r.path.Field("Roles").Map(role).Map(account)
	err := storage.Set[bool](rolePath, true)
	if err != nil {
		return
	}
}

func (r *RBAC) RevokeRole(role common.Hash, account common.Address) {
	rolePath := r.path.Field("Roles").Map(role).Map(account)
	err := storage.Set[bool](rolePath, false)
	if err != nil {
		return
	}
}

func (r *RBAC) SetupOwnerAsAdmin(owner common.Address) {
	r.GrantRole(RoleAdmin, owner)
}

func ComputeRoleHash(roleName string) common.Hash {
	return common.BytesToHash(common.LeftPadBytes([]byte(roleName), 32))
}
