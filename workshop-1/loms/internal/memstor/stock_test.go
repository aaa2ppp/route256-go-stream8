package memstor

import (
	"context"
	"route256/loms/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStock(t *testing.T) {
	ctx := context.Background()
	stock := NewStock()

	// Добавим тестовые данные
	stock.items[model.SKU(1)] = stockItem{count: 10, reserved: 0}
	stock.items[model.SKU(2)] = stockItem{count: 5, reserved: 0}
	// 	id	count	reserved
	// 	1	10		0
	//	2	5		0

	t.Run("GetInfo - existing SKU", func(t *testing.T) {
		count, err := stock.GetInfo(ctx, model.SKU(1))
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), count)
	})

	t.Run("GetInfo - non-existing SKU", func(t *testing.T) {
		_, err := stock.GetInfo(ctx, model.SKU(3))
		assert.ErrorIs(t, err, model.ErrNotFound)
	})

	t.Run("Reserve - success", func(t *testing.T) {
		err := stock.Reserve(ctx, []model.OrderItem{{SKU: model.SKU(1), Count: 3}})
		assert.NoError(t, err)
		count, _ := stock.GetInfo(ctx, model.SKU(1))
		assert.Equal(t, uint64(7), count)
	})
	// 	id	count	reserved
	// 	1	10		3
	//	2	5		0

	t.Run("Reserve - insufficient available", func(t *testing.T) {
		err := stock.Reserve(ctx, []model.OrderItem{{SKU: model.SKU(2), Count: 10}})
		assert.Error(t, err)
	})

	t.Run("ReserveCancel - success", func(t *testing.T) {
		err := stock.ReserveCancel(ctx, []model.OrderItem{{SKU: model.SKU(1), Count: 2}})
		assert.NoError(t, err)
		count, _ := stock.GetInfo(ctx, model.SKU(1))
		assert.Equal(t, uint64(9), count)
	})
	// 	id	count	reserved
	// 	1	10		1
	//	2	5		0

	t.Run("ReserveRemove - success", func(t *testing.T) {
		err := stock.ReserveRemove(ctx, []model.OrderItem{{SKU: model.SKU(1), Count: 1}})
		assert.NoError(t, err)
		count, _ := stock.GetInfo(ctx, model.SKU(1))
		assert.Equal(t, uint64(9), count)
	})
	// 	id	count	reserved
	// 	1	9		0
	//	2	5		0

	t.Run("ReserveCancel - insufficient reserved", func(t *testing.T) {
		err := stock.ReserveCancel(ctx, []model.OrderItem{{SKU: model.SKU(1), Count: 1}})
		assert.Error(t, err)
	})

	t.Run("ReserveRemove - insufficient reserved", func(t *testing.T) {
		err := stock.ReserveRemove(ctx, []model.OrderItem{{SKU: model.SKU(1), Count: 10}})
		assert.Error(t, err)
	})
}
