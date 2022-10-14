package casbin_service

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/sirupsen/logrus"
	"go_server/routers/casbin/DB"
)

//RbacDomain 支持多业务（domain）的casbin包
type RbacDomain struct {
	enforcer *casbin.Enforcer
	domain   string
}

func NewCasbin(dbConn string) (ret *RbacDomain, err error) {
	adapter := gormadapter.NewAdapterByDB(DB.MysqlTool())
	Enforcer, err := casbin.NewEnforcer("conf/keymatch.conf", adapter)
	if err != nil {
		logrus.Warn("casbin new enforcer err :", err)
		return
	}
	//注册超级管理员权限判断
	Enforcer.AddFunction("checkSuperAdmin", func(arguments ...interface{}) (i interface{}, err error) {
		username := arguments[0].(string)
		role := arguments[1].(string)
		//检查用户名的角色是否为超级管理员
		return Enforcer.HasRoleForUser(username, role)
	})

	//从数据库加载策略
	Enforcer.LoadPolicy()

	Enforcer.EnableLog(true)
	ret = &RbacDomain{
		enforcer: Enforcer,
		domain:   "",
	}
	return ret, nil
}

//AddRoleForUser 添加用户角色
func (r *RbacDomain) AddRoleForUser(username, role string) (bool, error) {
	return r.enforcer.AddRoleForUser(username, role)
}

//DeleteRoleForUser 删除用户角色
func (r *RbacDomain) DeleteRoleForUser(username, role string) (bool, error) {
	return r.enforcer.DeleteRoleForUser(username, role)
}

//GetUserRoles 获取用户角色
func (r *RbacDomain) GetUserRoles(username string) ([]string, error) {
	return r.enforcer.GetUsersForRole(username)
}

//GetRoleUsers 获取角色下的用户
func (r *RbacDomain) GetRoleUsers(role string) ([]string, error) {
	return r.enforcer.GetUsersForRole(role)
}

//AddRoleForUserMulti 批量添加用户角色
func (r *RbacDomain) AddRoleForUserMulti(username string, roles []string) (bool, error) {
	return r.enforcer.AddGroupingPolicy(r.formatMulti(username, roles))
}

//DeleteRoleForUserMulti 批量删除用户角色
func (r *RbacDomain) DeleteRoleForUserMulti(username string, roles []string) (bool, error) {
	return r.enforcer.RemoveGroupingPolicies(r.formatMulti(username, roles))
}

//UpdateRoleForUserMulti 批量更新用户角色
func (r *RbacDomain) UpdateRoleForUserMulti(username string, roles []string) (bool, error) {
	//删除旧角色
	_, err := r.enforcer.DeleteRolesForUser(username)
	if err != nil {
		return false, err
	}
	//添加新角色
	return r.enforcer.AddGroupingPolicies(r.formatMulti(username, roles))
}

//DeleteUser 删除用户：用户的权限一并删除
func (r *RbacDomain) DeleteUser(user string) (bool, error) {
	ok1, err := r.enforcer.RemoveFilteredGroupingPolicy(0, user, "")
	if err != nil {
		return ok1, err
	}
	ok2, err := r.enforcer.RemoveFilteredPolicy(0, user, "")
	if err != nil {
		return ok2, err
	}
	return ok1 || ok2, nil
}

//DeleteRole 删除角色：角色下的用户和权限一并删除
func (r *RbacDomain) DeleteRole(role string) (bool, error) {
	ok1, err := r.enforcer.RemoveFilteredGroupingPolicy(1, role)
	if err != nil {
		return ok1, err
	}
	ok2, err := r.enforcer.RemoveFilteredPolicy(1, role)
	if err != nil {
		return ok2, err
	}
	return ok1 || ok2, nil
}

//AddPermission 添加角色or用户权限
func (r *RbacDomain) AddPermission(username, permission string) (bool, error) {
	return r.enforcer.AddPermissionForUser(username, permission)
}

//AddPermissionMulti 批量添加角色or用户权限
func (r *RbacDomain) AddPermissionMulti(username string, permissions []string) (bool, error) {
	return r.enforcer.AddPolicies(r.formatMulti(username, permissions))
}

//DeletePermission 删除角色或用户权限
func (r *RbacDomain) DeletePermission(username, permission string) (bool, error) {
	return r.enforcer.DeletePermissionForUser(username, permission)
}

//DeletePermissionMulti 批量删除角色或用户权限
func (r *RbacDomain) DeletePermissionMulti(username string, permissions []string) (bool, error) {
	return r.enforcer.RemovePolicies(r.formatMulti(username, permissions))
}

//GetPermissionsForRole 获取角色权限
func (r *RbacDomain) GetPermissionsForRole(role string) []string {
	list := r.enforcer.GetFilteredNamedPolicy("p", 0, role, "")
	ret := make([]string, 0, len(list))
	for _, v := range list {
		ret = append(ret, v[1]) //todo 权限目前取第二位
	}
	return ret
}

//UpdatePermissionsForRoleMulti 批量更新角色权限
func (r *RbacDomain) UpdatePermissionsForRoleMulti(role string, polices []string) (bool, error) {
	//删除所有权限
	ok, err := r.enforcer.RemoveFilteredPolicy(0, role, "")
	if err != nil {
		return ok, err
	}
	//添加角色权限
	return r.enforcer.AddPolicies(r.formatMulti(role, polices))
}

//RemovePolice 删除指定权限
func (r *RbacDomain) RemovePolice(policy string) (bool, error) {
	return r.enforcer.RemoveFilteredPolicy(1, policy)
}

//HasPermission 用户或者角色是否具备权限
func (r *RbacDomain) HasPermission(user, permission string) (bool, error) {
	return r.enforcer.Enforce(user, permission)
}

//GetEnforcer 返回*casbin.Enforcer 用来执行casbin原生方法
func (r *RbacDomain) GetEnforcer() *casbin.Enforcer {
	return r.enforcer
}

func (r *RbacDomain) formatMulti(username string, polices []string) [][]string {
	policeArr := make([][]string, 0, len(polices))
	for _, v := range polices {
		policeArr = append(policeArr, []string{username, v})
	}
	return policeArr
}
