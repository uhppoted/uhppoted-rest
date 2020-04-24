package acl

import (
	"context"
	"log"
)

func warn(ctx context.Context, operation string, err error) {
	ctx.Value("log").(*log.Logger).Printf("WARN  %-20s %v\n", operation, err)
}
