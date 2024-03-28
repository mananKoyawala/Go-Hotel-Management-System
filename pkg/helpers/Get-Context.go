package helpers

import (
	"context"

	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
)

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), utils.TimeOut)
}
