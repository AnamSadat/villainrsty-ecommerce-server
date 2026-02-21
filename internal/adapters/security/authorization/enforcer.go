package casbin

import (
	"villainrsty-ecommerce-server/internal/core/shared/errors"

	casbin "github.com/casbin/casbin/v3"
)

func NewEnforcer(modelPath, policyPath string) (*casbin.Enforcer, error) {
	e, err := casbin.NewEnforcer(modelPath, policyPath)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to create casbin enforcer", err)
	}

	if err := e.LoadPolicy(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to load casbin policy", err)
	}

	return e, nil
}
