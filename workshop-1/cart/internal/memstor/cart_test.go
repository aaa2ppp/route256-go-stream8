package memstor

import (
	"context"
	"route256/cart/internal/model"
	"testing"
)

func TestCart_Add(t *testing.T) {
	cart := NewCart()

	tests := []struct {
		name      string
		userID    model.UserID
		req       model.AddCartItemRequest
		wantItems []model.CartItem
		wantErr   bool
	}{
		{
			name:   "add new item",
			userID: 1,
			req: model.AddCartItemRequest{
				UserID: 1,
				Items: []model.CartItem{
					{
						SKU:   123,
						Count: 2,
					},
				},
			},
			wantItems: []model.CartItem{{SKU: 123, Count: 2}},
			wantErr:   false,
		},
		{
			name:   "append new item",
			userID: 1,
			req: model.AddCartItemRequest{
				UserID: 1,
				Items: []model.CartItem{
					{
						SKU:   124,
						Count: 3,
					},
				},
			},
			wantItems: []model.CartItem{{SKU: 123, Count: 2}, {SKU: 124, Count: 3}},
			wantErr:   false,
		},
		{
			name:   "increase existing item count",
			userID: 1,
			req: model.AddCartItemRequest{
				UserID: 1,
				Items: []model.CartItem{
					{
						SKU:   123,
						Count: 3,
					},
				},
			},
			wantItems: []model.CartItem{{SKU: 123, Count: 5}, {SKU: 124, Count: 3}},
			wantErr:   false,
		},
		{
			name:   "add new two items",
			userID: 2,
			req: model.AddCartItemRequest{
				UserID: 2,
				Items: []model.CartItem{
					{
						SKU:   123,
						Count: 3,
					},
					{
						SKU:   124,
						Count: 4,
					},
				},
			},
			wantItems: []model.CartItem{{SKU: 123, Count: 3}, {SKU: 124, Count: 4}},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cart.Add(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			items, err := cart.List(context.Background(), tt.userID)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}

			if len(items) != len(tt.wantItems) {
				t.Errorf("expected %d items, got %d", len(tt.wantItems), len(items))
			}

			for i, item := range items {
				if item.SKU != tt.wantItems[i].SKU || item.Count != tt.wantItems[i].Count {
					t.Errorf("item %d: expected %v, got %v", i, tt.wantItems[i], item)
				}
			}
		})
	}
}

func TestCart_Delete(t *testing.T) {
	cart := NewCart()

	// Добавляем товар для тестирования удаления
	err := cart.Add(context.Background(), model.AddCartItemRequest{
		UserID: 1,
		Items: []model.CartItem{
			{
				SKU:   123,
				Count: 2,
			},
		},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	tests := []struct {
		name      string
		req       model.DeleteCartItemRequest
		wantItems []model.CartItem
		wantErr   bool
	}{
		{
			name: "delete existing item",
			req: model.DeleteCartItemRequest{
				UserID: 1,
				SKU:    123,
			},
			wantItems: []model.CartItem{},
			wantErr:   false,
		},
		{
			name: "delete non-existing item",
			req: model.DeleteCartItemRequest{
				UserID: 1,
				SKU:    456,
			},
			wantItems: []model.CartItem{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cart.Delete(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			items, err := cart.List(context.Background(), tt.req.UserID)
			if err != nil && !tt.wantErr {
				t.Fatalf("List() error = %v", err)
			}

			if len(items) != len(tt.wantItems) {
				t.Errorf("expected %d items, got %d", len(tt.wantItems), len(items))
			}
		})
	}
}

func TestCart_List(t *testing.T) {
	cart := NewCart()

	// Добавляем товар для тестирования
	err := cart.Add(context.Background(), model.AddCartItemRequest{
		UserID: 1,
		Items: []model.CartItem{
			{
				SKU:   123,
				Count: 2,
			},
		},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	tests := []struct {
		name      string
		userID    model.UserID
		wantItems []model.CartItem
		wantErr   bool
	}{
		{
			name:      "list items for existing user",
			userID:    1,
			wantItems: []model.CartItem{{SKU: 123, Count: 2}},
			wantErr:   false,
		},
		{
			name:      "list items for non-existing user",
			userID:    2,
			wantItems: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := cart.List(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(items) != len(tt.wantItems) {
				t.Errorf("expected %d items, got %d", len(tt.wantItems), len(items))
			}

			for i, item := range items {
				if item.SKU != tt.wantItems[i].SKU || item.Count != tt.wantItems[i].Count {
					t.Errorf("item %d: expected %v, got %v", i, tt.wantItems[i], item)
				}
			}
		})
	}
}

func TestCart_Clear(t *testing.T) {
	cart := NewCart()

	// Добавляем товар для тестирования
	err := cart.Add(context.Background(), model.AddCartItemRequest{
		UserID: 1,
		Items: []model.CartItem{
			{
				SKU:   123,
				Count: 2,
			},
		},
	})
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	tests := []struct {
		name    string
		userID  model.UserID
		wantErr bool
	}{
		{
			name:    "clear existing user cart",
			userID:  1,
			wantErr: false,
		},
		{
			name:    "clear non-existing user cart",
			userID:  2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cart.Clear(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Clear() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			items, err := cart.List(context.Background(), tt.userID)
			if err == nil {
				t.Errorf("expected cart to be empty, got %v", items)
			}
		})
	}
}
