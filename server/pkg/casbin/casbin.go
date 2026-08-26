package casbin

import (
	"interastral-peace.com/alnitak/pkg/mysql"
	"interastral-peace.com/alnitak/utils"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
)

type Casbin struct {
	casbinEnforcer *casbin.Enforcer
}

func InitCasbin() *Casbin {
	a, err := gormadapter.NewAdapterByDB(mysql.GetMysqlClient())
	if err != nil {
		utils.ErrorLog("casbin初始化失败", "casbin", err.Error())
		return nil
	}

	e, err := casbin.NewEnforcer("./static/casbin/model.conf", a)
	if err != nil {
		utils.ErrorLog("casbin初始化失败", "casbin", err.Error())
		return nil
	}

	if err := e.LoadPolicy(); err != nil {
		utils.ErrorLog("casbin初始化失败", "casbin", err.Error())
		return nil
	}

	zap.L().Info("casbin初始化成功", zap.String("module", "casbin"))

	return &Casbin{
		casbinEnforcer: e,
	}

}

func (c *Casbin) CasbinCheck(sub string, obj string, act string) bool {
	// Enforce 内部已使用 RWMutex，并发安全，无需额外加锁
	pass, err := c.casbinEnforcer.Enforce(sub, obj, act)
	if err != nil {
		utils.ErrorLog("casbin校验失败", "casbin", err.Error())
		return false
	}
	return pass
}

func (c *Casbin) DeletePolicy(sub string, obj string, act string) bool {
	ok, err := c.casbinEnforcer.RemovePolicy(sub, obj, act)
	if !ok {
		utils.ErrorLog("移除casbin policy失败", "casbin", err.Error())
	}
	return ok
}

func (c *Casbin) AddPolicy(sub string, obj string, act string) bool {
	ok, err := c.casbinEnforcer.AddPolicy(sub, obj, act)
	if !ok {
		utils.ErrorLog("添加casbin policy失败", "casbin", err.Error())
	}
	return ok
}

// ReloadPolicy 重新从数据库加载策略
func (c *Casbin) ReloadPolicy() error {
	return c.casbinEnforcer.LoadPolicy()
}
