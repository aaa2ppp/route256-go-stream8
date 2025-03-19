// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package cart

import (
	"context"
)

type Querier interface {
	AddArrays(ctx context.Context, arg AddArraysParams) error
	AddBatch(ctx context.Context, arg []AddBatchParams) *AddBatchBatchResults
	Clear(ctx context.Context, userID int64) error
	Delete(ctx context.Context, arg DeleteParams) error
	List(ctx context.Context, userID int64) ([]Cart, error)
}

var _ Querier = (*Queries)(nil)
