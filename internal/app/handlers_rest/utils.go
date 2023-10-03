package handlersrest

import (
	"context"
	"errors"
	"fmt"
)

const defaultServiceError = "service error:"

func getUserID(ctx context.Context) (string, error) {
	value := ctx.Value(userCtx)
	if value == nil {
		return "", errors.New("nil token")
	}
	id, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("user id is of invalid type %T", value)
	}
	return id, nil
}
