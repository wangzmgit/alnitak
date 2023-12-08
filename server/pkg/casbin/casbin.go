package casbin

import (
	"fmt"
	"sync"

	"interastral-peace.com/alnitak/pkg/mysql"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
)

var checkLock sync.Mutex

type Casbin struct {
	casbinEnforcer *casbin.Enforcer
}

func InitCasbin() *Casbin {
	a, err := gormadapter.NewAdapterByDB(mysql.GetMysqlClient())
	if err != nil {
		zap.L().Error(fmt.Sprintf("初始化错误 | Casbin | 原因: %s", err.Error()))
		return nil
	}

	e, err := casbin.NewEnforcer("./static/casbin/model.conf", a)
	if err != nil {
		zap.L().Error(fmt.Sprintf("初始化错误 | Casbin | 原因: %s", err.Error()))

		return nil
	}

	if err := e.LoadPolicy(); err != nil {
		zap.L().Error(fmt.Sprintf("初始化错误 | Casbin | 原因: %s", err.Error()))

		return nil
	}

	zap.L().Info("初始化成功 | Casbin | casbin初始化成功")

	return &Casbin{
		casbinEnforcer: e,
	}

}

func (c *Casbin) CasbinCheck(sub string, obj string, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()

	pass, _ := c.casbinEnforcer.Enforce(sub, obj, act)

	return pass
}

func (c *Casbin) DeletePolicy(sub string, obj string, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()

	ok, err := c.casbinEnforcer.RemovePolicy(sub, obj, act)
	if !ok {
		zap.L().Error(fmt.Sprintf("Casbin错误 | 移除casbin policy 失败 | 原因: %s", err.Error()))
	}
	return ok
}

func (c *Casbin) AddPolicy(sub string, obj string, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()

	ok, err := c.casbinEnforcer.AddPolicy(sub, obj, act)
	if !ok {
		zap.L().Error(fmt.Sprintf("Casbin错误 | 添加casbin policy 失败 | 原因: %s", err.Error()))
	}
	return ok
}
