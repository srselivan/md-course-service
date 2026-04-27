package utils

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func GetInt64Param(ctx fiber.Ctx, key string) (int64, error) {
	param := ctx.Params(key)
	if param == "" {
		return 0, fmt.Errorf("param %s is not exist", key)
	}
	paramInt64, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt(%s) fail: %w", param, err)
	}
	return paramInt64, nil
}

func GetInt64Query(ctx fiber.Ctx, key string) (int64, error) {
	query := ctx.Query(key)
	if query == "" {
		return 0, fmt.Errorf("query %s is not exist", key)
	}
	queryInt64, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt(%s) fail: %w", query, err)
	}
	return queryInt64, nil
}

func GetInt64pQuery(ctx fiber.Ctx, key string) (*int64, error) {
	query := ctx.Query(key)
	if query == "" {
		return nil, nil
	}
	queryInt64, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("strconv.ParseInt(%s) fail: %w", query, err)
	}
	return &queryInt64, nil
}

func Int64pToInt16p(src *int64) *int16 {
	if src == nil {
		return nil
	}
	return new(int16(*src))
}
