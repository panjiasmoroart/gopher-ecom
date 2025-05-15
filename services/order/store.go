package order

import (
	"database/sql"

	"github.com/panjiasmoroart/gopher-ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateOrderWithTx - Membuat pesanan menggunakan transaksi
func (s *Store) CreateOrder(tx *sql.Tx, order types.Order) (int, error) {
	result, err := tx.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(orderID), nil
}

// CreateOrderItemWithTx - Membuat order item menggunakan transaksi
func (s *Store) CreateOrderItem(tx *sql.Tx, orderItem types.OrderItem) error {
	_, err := tx.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}

func (s *Store) CreateOrderWithoutTx(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItemWithoutTx(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}
