package order

import (
	"context"
	"reflect"
	"route256/loms/internal/model"
	"route256/loms/internal/repo/order/queries"
	"route256/loms/internal/service"

	"github.com/jackc/pgx/v5"
)

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	queries.DBTX
}

type Adapter struct {
	db DB
	q  *queries.Queries
}

func New(db DB) *Adapter {
	return &Adapter{
		db: db,
		q:  queries.New(db),
	}
}

// Create implements storage.
func (a Adapter) Create(ctx context.Context, req model.CreateOrderRequest) (model.OrderID, error) {
	const op = "Create"
	log := getLogger(ctx, op)

	var orderID model.OrderID

	err := pgx.BeginFunc(ctx, a.db, func(tx pgx.Tx) error {
		qtx := a.q.WithTx(tx)

		if id, err := qtx.Create(ctx, int64(req.UserID)); err != nil {
			return err
		} else {
			orderID = model.OrderID(id)
		}

		n := len(req.Items)
		orders := make([]int64, 0, n)
		skus := make([]int64, 0, n)
		counts := make([]int32, 0, n)
		for _, item := range req.Items {
			orders = append(orders, int64(orderID))
			skus = append(skus, int64(item.SKU))
			counts = append(counts, int32(item.Count))
		}

		if err := qtx.AddItems(ctx, queries.AddItemsParams{
			Orders: orders,
			Skus:   skus,
			Counts: counts,
		}); err != nil {
			return err
		}

		log.Debug("crete order done", "orderID", orderID)
		return nil
	})

	if err != nil {
		log.Debug("BeginFunc", "error", err, "errorType", reflect.TypeOf(err))
		return 0, err
	}

	return orderID, nil
}

// GetByID implements storage.
func (a Adapter) GetByID(ctx context.Context, orderID model.OrderID) (model.Order, error) {
	qresp, err := a.q.GetByID(ctx, int64(orderID))
	if err != nil {
		return model.Order{}, err
	}
	if len(qresp) == 0 {
		return model.Order{}, model.ErrNotFound
	}
	status, err := model.ParseOrderStatus(string(qresp[0].Status))
	if err != nil {
		return model.Order{}, err
	}

	resp := model.Order{
		UserID: model.UserID(qresp[0].UserID),
		Status: status,
		Items:  make([]model.OrderItem, 0, len(qresp)),
	}

	for _, item := range qresp {
		resp.Items = append(resp.Items, model.OrderItem{
			SKU:   model.SKU(item.Sku),
			Count: uint16(item.Count),
		})
	}

	return resp, nil
}

// SetStatus implements storage.
func (a Adapter) SetStatus(ctx context.Context, req model.SetOrderStatusRequest) error {
	return a.q.SetStatus(ctx, queries.SetStatusParams{
		OrderID: int64(req.OrderID),
		Status:  queries.OrderStatus(req.Status.String()),
	})
}

var _ service.OrderStorage = Adapter{}
