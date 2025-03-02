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
	GetInfo(ctx context.Context, req model.GetProductRequest) (model.GetProductResponse, error)
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
	countInCart, err := p.getCountInCart(ctx, req.UserID)
	if err != nil {
		return err
	}
	sreq := model.AddCartItemRequest{
		UserID: req.UserID,
		Items:  make([]model.CartItem, 0, len(req.Items)),
	}
	for _, item := range req.Items {
		if item.Count == 0 {
			// Игнорируем нулевые количества
			continue
		}

		exists, err := p.isProductExists(ctx, req.Token, item.SKU)
		if err != nil {
			return err
		}
		if !exists {
			return model.ErrPreconditionFailed
		}

		count, err := p.getCountInStock(ctx, item.SKU)
		if err != nil {
			return err
		}
		if count < uint64(item.Count+countInCart[item.SKU]) {
			return model.ErrPreconditionFailed
		}

		sreq.Items = append(sreq.Items, item)
	}
	return p.cart.Add(ctx, sreq)
}

func (p *Cart) getCountInCart(ctx context.Context, userID model.UserID) (map[model.SKU]uint16, error) {
	cartItems, err := p.cart.List(ctx, userID)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("failed to list cart items: %w", err)
	}
	countInCart := make(map[model.SKU]uint16, len(cartItems))
	for _, item := range cartItems {
		countInCart[item.SKU] = item.Count
	}
	return countInCart, nil
}

func (p *Cart) isProductExists(ctx context.Context, token string, sku model.SKU) (bool, error) {
	_, err := p.product.GetInfo(ctx, model.GetProductRequest{
		Token: token,
		SKU:   sku,
	})
	if errors.Is(err, model.ErrNotFound) {
		// Товар не найден в каталоге, но здесь это не ошибка
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to get product info: %w", err)
	}
	return true, nil
}

func (p *Cart) getCountInStock(ctx context.Context, sku model.SKU) (count uint64, _ error) {
	count, err := p.order.GetStockInfo(ctx, sku)
	if errors.Is(err, model.ErrNotFound) {
		// Если товар не найден на складе, считаем его количество равным 0
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get product in stock info: %w", err)
	}
	return count, nil
}

func (p *Cart) Delete(ctx context.Context, req model.DeleteCartItemRequest) error {
	err := p.cart.Delete(ctx, req)
	if errors.Is(err, model.ErrNotFound) {
		// Если товара нет в корзине, не возвращаем ошибку, чтобы обеспечить идемпотентность
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to delete item from cart: %w", err)
	}
	return nil
}

func (p *Cart) List(ctx context.Context, req model.CartListRequest) (resp model.CartListResponse, _ error) {
	log := getLogger(ctx, "Cart.List")

	// Инициализируем список товаров в корзине, чтобы не вернуть nil
	resp.Items = []model.CartListItem{}

	cartItems, err := p.cart.List(ctx, req.UserID)
	if errors.Is(err, model.ErrNotFound) {
		return resp, nil
	}
	if err != nil {
		return resp, fmt.Errorf("failed to list cart items: %w", err)
	}

	var skusToDelete []model.SKU
	for _, item := range cartItems {
		if item.Count == 0 {
			// Товара нет в корзине
			log.Warn("no product in cart, it will be deleted", "SKU", item.SKU)
			skusToDelete = append(skusToDelete, item.SKU)
			continue
		}
		product, err := p.product.GetInfo(ctx, model.GetProductRequest{
			Token: req.Token,
			SKU:   item.SKU,
		})
		if errors.Is(err, model.ErrNotFound) {
			// Товар не найден в каталоге, добавляем в список на удаление
			log.Warn("cart item not found in products, it will be deleted", "SKU", item.SKU)
			skusToDelete = append(skusToDelete, item.SKU)
			continue
		}
		if err != nil {
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
	for _, sku := range skusToDelete {
		if err := p.cart.Delete(ctx, model.DeleteCartItemRequest{
			UserID: req.UserID,
			SKU:    sku,
		}); err != nil {
			return resp, fmt.Errorf("failed to delete invalid cart item: %w", err)
		}
	}

	return resp, nil
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
