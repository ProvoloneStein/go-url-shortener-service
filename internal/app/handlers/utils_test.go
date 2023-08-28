package handlers

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_getUserID(t *testing.T) {
	// Init Test Table
	type want struct {
		res string
		err error
	}

	tests := []struct {
		name string
		ctx  context.Context
		want want
	}{
		{
			name: "Good test",
			ctx:  context.WithValue(context.Background(), userCtx, "3213"),
			want: want{
				res: "3213",
				err: nil,
			},
		},
		{
			name: "err test",
			ctx:  context.WithValue(context.Background(), userCtx, 12),
			want: want{
				res: "",
				err: errors.New("user id is of invalid type int"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Create Request
			res, err := getUserID(tt.ctx)
			assert.Equal(t, tt.want.res, res)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
