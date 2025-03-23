package service

import (
	"context"
	"errors"
	"fmt"
	"route256/cart/internal/model"
)

type CartStorage interface {
	Add(ctx context.Context, req model.AddCartItemRequest) error
	Delete(ctx context.Context, req model.DeleteCartItemRequest) error
	List(ctx context.Context, userID model.UserID) ([]model.CartItem, error)
	Clear(ctx context.Context, userID model.UserID) error
}

type OrderStorage interface {
	CreateOrder(ctx context.Context, req model.OrderCreateRequest) (model.OrderID, error)
	GetStockInfo(ctx context.Context, sku model.SKU) (count uint64, err error)
}

type ProductStorage interface {
	GetInfo(ctx context.Context, sku model.SKU) (model.GetProductResponse, error)
}

type Cart struct {
	cart    CartStorage
	order   OrderStorage
	product ProductStorage
}

func NewCart(cart CartStorage, order OrderStorage, product ProductStorage) *Cart {
	return &Cart{
		cart:    cart,
		order:   order,
		product: product,
	}
}

func (p *Cart) Add(ctx context.Context, req model.AddCartItemRequest) error {
	if err := p.checkProduct(ctx, req.SKU); err != nil {
		return err
	}

	count, err := p.getCountInStock(ctx, req.SKU)
	if err != nil {
		return err
	}

	if count < uint64(req.Count) {
		return model.ErrPreconditionFailed
	}

	return p.cart.Add(ctx, req)
}

func (p *Cart) checkProduct(ctx context.Context, sku model.SKU) error {
	if _, err := p.product.GetInfo(ctx, sku); err != nil {
		return fmt.Errorf("failed to get product info: %w", err)
	}
	return nil
}

func (p *Cart) getCountInStock(ctx context.Context, sku model.SKU) (count uint64, _ error) {
	count, err := p.order.GetStockInfo(ctx, sku)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// Если товар не найден на складе, считаем его количество равным 0
			return 0, nil
		}
		return 0, fmt.Errorf("failed to get product in stock info: %w", err)
	}
	return count, nil
}

func (p *Cart) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	err := p.cart.Delete(ctx, req)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// Если товара нет в корзине, не возвращаем ошибку, чтобы обеспечить идемпотентность
			return nil
		}
		return fmt.Errorf("failed to delete item from cart: %w", err)
	}
	return nil
}

func (p *Cart) List(ctx context.Context, userID model.UserID) (resp model.CartListResponse, _ error) {
	log := getLogger(ctx, "Cart.List")

	// Инициализируем список товаров в корзине, чтобы не вернуть nil
	resp.Items = []model.CartListItem{}

	cartItems, err := p.cart.List(ctx, userID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			// Если не найдено ни одного товара в корзине, считаем что корзина пуста
			return resp, nil
		}
		return resp, fmt.Errorf("failed to list cart items: %w", err)
	}

	var nonExistents []model.DeleteCartItemRequest
	for _, item := range cartItems {
		product, err := p.product.GetInfo(ctx, item.SKU)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				// Товар не найден в каталоге, добавляем в список на удаление
				log.Warn("cart item not found in products, it will be deleted", "SKU", item.SKU)
				nonExistents = append(nonExistents, model.DeleteCartItemRequest{
					UserID: userID,
					SKU:    item.SKU,
				})
				continue
			}
			return resp, fmt.Errorf("failed to get product info: %w", err)
		}

		// Добавляем товар в ответ
		resp.Items = append(resp.Items, model.CartListItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  product.Name,
			Price: product.Price,
		})
		resp.TotalPrice += product.Price * uint32(item.Count)
	}

	// Удаляем несуществующие товары из корзины
	if err := p.deleteBatch(ctx, nonExistents); err != nil {
		log.Error("failed to delete invalid cart item", "error", err)
	}

	return resp, nil
}

func (p *Cart) deleteBatch(ctx context.Context, batch []model.DeleteCartItemRequest) error {
	// TODO: нужен бачевый метод в репо?
	for _, item := range batch {
		if err := p.cart.Delete(ctx, item); err != nil {
			// что-то совсем пошло не так
			return fmt.Errorf(": %w", err)
		}
	}
	return nil
}

func (p *Cart) Clear(ctx context.Context, userID model.UserID) error {
	return p.cart.Clear(ctx, userID)
}

func (p *Cart) Checkout(ctx context.Context, userID model.UserID) (model.OrderID, error) {
	log := getLogger(ctx, "Cart.Checkout")

	cartItems, err := p.cart.List(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to list cart items: %w", err)
	}

	if len(cartItems) == 0 {
		return 0, fmt.Errorf("cart is empty: %w", model.ErrNotFound)
	}

	orderID, err := p.order.CreateOrder(ctx, model.OrderCreateRequest{
		UserID: userID,
		Items:  cartItems,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	if err := p.cart.Clear(ctx, userID); err != nil {
		// Логируем ошибку очистки корзины, но не прерываем выполнение
		log.Error("failed to clear cart after checkout", "error", err)
	}

	return orderID, nil
}
