package handlers

import (
	"context"
	"errors"
)

func getUserId(ctx context.Context) (string, error) {
	id, ok := ctx.Value(userCtx).(string)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}
	return id, nil
}
