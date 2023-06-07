package permission

import (
	"context"
	"errors"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"github.com/oarkflow/frame"
	"github.com/oarkflow/frame/pkg/protocol/consts"
	gormadapter "github.com/oarkflow/gorm-adapter"
	"github.com/oarkflow/govaluate"
	"gorm.io/gorm"

	"github.com/oarkflow/pkg/str"
)

var Instance *Engine

type Config struct {
	TableName       string
	PrimaryKey      string
	DisableMigrate  bool
	Unauthorized    frame.HandlerFunc
	Forbidden       frame.HandlerFunc
	ParamExtractor  func(c context.Context, ctx *frame.Context) []string
	CustomFunctions map[string]govaluate.ExpressionFunction
	DB              *gorm.DB
	Model           interface{}
	Policy          interface{}
	Adapter         persist.Adapter
	Scopes          []func(db *gorm.DB) *gorm.DB
}

// Engine holds the configuration for the middleware
type Engine struct {
	*Enforcer
	PolicyAdapter persist.Adapter
	config        Config
}

func Default(cfg Config) (*Engine, error) {
	engine, err := New(cfg)
	Instance = engine
	return engine, err
}

func New(cfg Config) (*Engine, error) {
	var params []any
	if cfg.Model == nil {
		cfg.Model = roleModel()
	}
	params = append(params, cfg.Model)
	if cfg.Policy != nil {
		params = append(params, cfg.Policy)
	}
	if cfg.TableName == "" {
		cfg.TableName = "permissions"
	}
	if cfg.ParamExtractor == nil {
		return nil, errors.New("Parameter extractor not defined")
	}

	if cfg.Unauthorized == nil {
		cfg.Unauthorized = func(cc context.Context, c *frame.Context) {
			c.AbortWithJSON(consts.StatusUnauthorized, consts.StatusUnauthorized)
			return
		}
	}
	if cfg.Forbidden == nil {
		cfg.Forbidden = func(cc context.Context, c *frame.Context) {
			c.AbortWithJSON(consts.StatusForbidden, consts.StatusForbidden)
			return
		}
	}

	enforcer, err := casbin.NewEnforcer(params...)
	if err != nil {
		return nil, err
	}
	engine := &Engine{
		Enforcer: &Enforcer{Enforcer: enforcer},
		config:   cfg,
	}
	if len(funcs) > 0 {
		for key, fn := range funcs {
			engine.AddFunction(key, fn)
		}
	}
	if len(cfg.CustomFunctions) > 0 {
		for key, fn := range cfg.CustomFunctions {
			engine.AddFunction(key, fn)
		}
	}
	if cfg.DB != nil {
		if cfg.DisableMigrate {
			gormadapter.TurnOffAutoMigrate(cfg.DB)
		}
		adapter, err := gormadapter.New(gormadapter.Config{
			TableName:  cfg.TableName,
			PrimaryKey: cfg.PrimaryKey,
			DB:         cfg.DB,
			Scopes:     cfg.Scopes,
		})
		if err != nil {
			return nil, err
		}
		engine.PolicyAdapter = adapter
		enforcer.SetAdapter(adapter)
	}
	if cfg.Adapter != nil && engine.PolicyAdapter == nil {
		engine.PolicyAdapter = cfg.Adapter
		enforcer.SetAdapter(cfg.Adapter)
	}
	enforcer.EnableAutoSave(true)
	err = enforcer.LoadPolicy()
	return engine, err
}

// RequirePermissions tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
func (cm *Engine) RequirePermissions(permissions []string, opts ...func(o *Options)) frame.HandlerFunc {
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(cc context.Context, c *frame.Context) {
		if len(permissions) == 0 {
			c.Next(cc)
			return
		}
		switch options.ValidationRule {
		case matchAll:
			for _, permission := range permissions {
				vals := options.PermissionParser(permission)
				if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
					c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
					return
				} else if !ok {
					cm.config.Forbidden(cc, c)
					return
				}
			}
			c.Next(cc)
			return
		case atLeastOne:
			for _, permission := range permissions {
				vals := options.PermissionParser(permission)
				if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
					c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
					return
				} else if ok {
					c.Next(cc)
					return
				}
			}
			cm.config.Forbidden(cc, c)
			return
		}
		c.Next(cc)
		return
	}
}

// Can try to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
func (cm *Engine) Can(dom, sub, perm string, opts ...func(o *Options)) bool {
	permissions := []string{perm}
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}
	if len(permissions) == 0 {
		return false
	}
	switch options.ValidationRule {
	case matchAll:
		for _, permission := range permissions {
			vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
			if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
				return false
			} else if !ok {
				return false
			}
		}
		return true
	case atLeastOne:
		for _, permission := range permissions {
			vals := append([]string{sub, dom}, options.PermissionParser(permission)...)
			if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
				return false
			} else if ok {
				return true
			}
		}
		return false
	}
	return false
}

// RoutePermission tries to find the current subject and determine if the
// subject has the required permissions according to predefined Casbin policies.
// This method uses http Path and Method as object and action.
func (cm *Engine) RoutePermission(cc context.Context, c *frame.Context) {
	vals := cm.config.ParamExtractor(cc, c)
	if len(vals) < 3 {
		cm.config.Unauthorized(cc, c)
		return
	}
	sub := vals[0]
	if sub == "" {
		cm.config.Unauthorized(cc, c)
		return
	}
	dom := vals[1]
	if dom == "" {
		dom = "*"
	}
	availableDomains, _ := cm.Enforcer.GetDomainsForUser(sub)
	availableDomainCount := len(availableDomains)
	if !str.Contains(availableDomains, dom) && availableDomainCount > 1 {
		cm.config.Forbidden(cc, c)
		return
	} else if availableDomainCount == 1 && availableDomains[0] != "*" && availableDomains[0] != dom {
		cm.config.Forbidden(cc, c)
		return
	}

	if str.Contains(availableDomains, "*") && dom != "*" {
		dom = "*"
	}
	if ok, err := cm.Enforcer.Enforce(str.ConvertToInterface(vals)...); err != nil {
		c.AbortWithJSON(consts.StatusInternalServerError, err.Error())
		return
	} else if !ok {
		cm.config.Forbidden(cc, c)
		return
	}

	c.Next(cc)
	return
}

// RequireRoles tries to find the current subject and determine if the
// subject has the required roles according to predefined Casbin policies.
func (cm *Engine) RequireRoles(roles []string, opts ...func(o *Options)) frame.HandlerFunc {
	options := &Options{
		ValidationRule:   matchAll,
		PermissionParser: permissionParserWithSeparator(":"),
	}

	for _, o := range opts {
		o(options)
	}

	return func(cc context.Context, c *frame.Context) {
		if len(roles) == 0 {
			c.Next(cc)
			return
		}
		userRoles := []string{}
		if options.ValidationRule == matchAll {
			for _, role := range roles {
				if !str.Contains(userRoles, role) {
					cm.config.Forbidden(cc, c)
					return
				}
			}
			c.Next(cc)
			return
		} else if options.ValidationRule == atLeastOne {
			for _, role := range roles {
				if str.Contains(userRoles, role) {
					c.Next(cc)
					return
				}
			}
			cm.config.Forbidden(cc, c)
			return
		}

		c.Next(cc)
		return
	}
}
