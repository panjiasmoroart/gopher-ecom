package cart

import (
	"fmt"
	"log"

	"github.com/panjiasmoroart/gopher-ecom/types"
)

func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	productIds := make([]int, len(items))
	for index, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[index] = item.ProductID
	}

	return productIds, nil
}

func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

// # createOrder mengembalikan (int, float64, error):
// int: Total kuantitas produk dalam pesanan.
// float64: Total harga pesanan.
// error: Error (jika ada), terutama dari pengecekan stok.

func (h *Handler) createOrder(products []types.Product, cartItems []types.CartCheckoutItem, userID int) (int, float64, error) {
	// Create a map of products for easier access
	productsMap := make(map[int]types.Product)
	for _, product := range products {
		productsMap[product.ID] = product
	}

	// map[1:{1 Apple.jpg 1.5 100 2025-05-16 14:00:00 +0000 UTC}
	//  2:{2 Banana.jpg 1.0 200 2025-05-16 14:00:00 +0000 UTC}]

	// Check if all products are available
	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
		return 0, 0, err
	}

	// Calculate total price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	log.Println("starting transaction db")
	// start transaction database
	tx, err := h.db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// reduce/kurangi the quantity of products in the store
	for _, item := range cartItems {
		product := productsMap[item.ProductID]
		product.Quantity -= item.Quantity
		h.store.UpdateProduct(tx, product)
	}

	// Create order record
	orderID, err := h.orderStore.CreateOrder(tx, types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "Jl. Pemuda Sawangan Baru - Depok",
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create order: %v", err)
	}

	// Create order items and update product stock
	for _, item := range cartItems {
		// Create Order Item
		err = h.orderStore.CreateOrderItem(tx, types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to create order item: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return orderID, totalPrice, nil
}

// func (h *Handler) createOrderOld(products []types.Product, cartItems []types.CartCheckoutItem) (int, float64, error) {
// 	// create a map of products for easier access
// 	productsMap := make(map[int]types.Product)
// 	for _, product := range products {
// 		productsMap[product.ID] = product
// 	}

// 	// check if all products are available
// 	if err := checkIfCartIsInStock(cartItems, productsMap); err != nil {
// 		return 0, 0, err
// 	}

// 	// calculate total price
// 	totalPrice := calculateTotalPrice(cartItems, productsMap)

// 	// reduce/kurangi the quantity of products in the store
// 	for _, item := range cartItems {
// 		product := productsMap[item.ProductID]
// 		product.Quantity -= item.Quantity
// 		h.store.UpdateProduct(product)
// 	}

// 	// create order recored
// 	orderID, err := h.orderStore.CreateOrder(types.Order{
// 		UserID:  1,
// 		Total:   totalPrice,
// 		Status:  "pending",
// 		Address: "Jl. Pemuda Sawangan Baru - Depok",
// 	})
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	// create order the items records
// 	for _, item := range cartItems {
// 		h.orderStore.CreateOrderItem(types.OrderItem{
// 			OrderID:   orderID,
// 			ProductID: item.ProductID,
// 			Quantity:  item.Quantity,
// 			Price:     productsMap[item.ProductID].Price,
// 		})
// 	}

// 	return orderID, totalPrice, nil
// }
