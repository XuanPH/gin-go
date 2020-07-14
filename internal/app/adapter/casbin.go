package adapter

import (
	"context"
	"fmt"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/chunganhbk/gin-go/internal/app/repositories"
	"github.com/chunganhbk/gin-go/internal/app/schema"
	"github.com/chunganhbk/gin-go/pkg/logger"
)
/*type Adapter interface {
    // LoadPolicy loads all policy rules from the storage.
    LoadPolicy(models models.Model) error
    // SavePolicy saves all policy rules to the storage.
    SavePolicy(models models.Model) error
    // AddPolicy adds a policy rule to the storage.
    // This is part of the Auto-Save feature.
    AddPolicy(sec string, ptype string, rule []string) error
    // RemovePolicy removes a policy rule from the storage.
    // This is part of the Auto-Save feature.
    RemovePolicy(sec string, ptype string, rule []string) error
    // RemoveFilteredPolicy removes policy rules that match the filter from the storage.
    // This is part of the Auto-Save feature.
    RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
}*/


// Casbin Adapter casbin
type CasbinAdapter struct {
	RoleRp     repositories.IRole
	RoleMenuRp    repositories.IRoleMenu
	MenuResourceRp repositories.IMenuActionResource
	UserRp        repositories.IUser
	UserRoleRp     repositories.IUserRole
}
func NewCasbinAdapter(roleRp repositories.IRole, roleMenuRp repositories.IRoleMenu,
	menuResourceRp repositories.IMenuActionResource,
    userRp repositories.IUser, userRoleRp repositories.IUserRole) *CasbinAdapter {
	return &CasbinAdapter{
		RoleRp:   roleRp,
		RoleMenuRp:   roleMenuRp,
		MenuResourceRp: menuResourceRp,
		UserRp:    userRp,
		UserRoleRp:   userRoleRp,
	}
}
// LoadPolicy loads all policy rules from the storage.
func (c *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := c.loadRolePolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin role policy error: %s", err.Error())
		return err
	}

	err = c.loadUserPolicy(ctx, model)
	if err != nil {
		logger.Errorf(ctx, "Load casbin user policy error: %s", err.Error())
		return err
	}

	return nil
}

// load Role Policy  (p,role_id,path,method)
func (c *CasbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	roleResult, err := c.RoleRp.Query(ctx, schema.RoleQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return nil
	}

	roleMenuResult, err := c.RoleMenuRp.Query(ctx, schema.RoleMenuQueryParam{})
	if err != nil {
		return err
	}
	mRoleMenus := roleMenuResult.Data.ToRoleIDMap()

	menuResourceResult, err := c.MenuResourceRp.Query(ctx, schema.MenuActionResourceQueryParam{})
	if err != nil {
		return err
	}
	mMenuResources := menuResourceResult.Data.ToActionIDMap()

	for _, item := range roleResult.Data {
		mcache := make(map[string]struct{})
		if rms, ok := mRoleMenus[item.ID]; ok {
			for _, actionID := range rms.ToActionIDs() {
				if mrs, ok := mMenuResources[actionID]; ok {
					for _, mr := range mrs {
						if mr.Path == "" || mr.Method == "" {
							continue
						} else if _, ok := mcache[mr.Path+mr.Method]; ok {
							continue
						}
						mcache[mr.Path+mr.Method] = struct{}{}
						line := fmt.Sprintf("p,%s,%s,%s", item.ID, mr.Path, mr.Method)
						persist.LoadPolicyLine(line, m)
					}
				}
			}
		}
	}

	return nil
}

// load User Policy(g,user_id,role_id)
func (c *CasbinAdapter) loadUserPolicy(ctx context.Context, m casbinModel.Model) error {
	userResult, err := c.UserRp.Query(ctx, schema.UserQueryParam{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(userResult.Data) > 0 {
		userRoleResult, err := c.UserRoleRp.Query(ctx, schema.UserRoleQueryParam{})
		if err != nil {
			return err
		}

		mUserRoles := userRoleResult.Data.ToUserIDMap()
		for _, uitem := range userResult.Data {
			if urs, ok := mUserRoles[uitem.ID]; ok {
				for _, ur := range urs {
					line := fmt.Sprintf("g,%s,%s", ur.UserID, ur.RoleID)
					persist.LoadPolicyLine(line, m)
				}
			}
		}
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
