package repoimpl

import (
	"database/sql"
	"order-and-pay/intrface"
	"order-and-pay/model"
	"time"
)

type OrderRepoImpl struct {
	db *sql.DB
	q  intrface.Querier
}

func NewOrderRepoImpl(db *sql.DB) *OrderRepoImpl {
	return &OrderRepoImpl{db, db}
}

func (ori *OrderRepoImpl) GetPaymentsSumByOrderId(orderId string) (float32, error) {
	payRes, err := ori.q.Query("SELECT SUM(total)-SUM(change) FROM payments WHERE order_id=$1", orderId)
	if err != nil {
		return 0, err
	}
	defer payRes.Close()
	var paymentSum float32
	if payRes.Next() {
		payRes.Scan(&paymentSum)
	}
	return paymentSum, nil
}

func (ori *OrderRepoImpl) GetProductsPriceSumByOrderId(orderId string) (float32, error) {
	prodPriceSumRes, err := ori.q.Query("SELECT SUM(num*price_per_one) FROM products_in_orders WHERE order_id=$1", orderId)
	if err != nil {
		return 0, err
	}
	defer prodPriceSumRes.Close()
	var priceSum float32
	if prodPriceSumRes.Next() {
		prodPriceSumRes.Scan(&priceSum)
	}
	return priceSum, nil
}

func (ori *OrderRepoImpl) UpdateOrderStatusToComplete(orderId string) error {
	_, err := ori.q.Query("UPDATE orders SET status_id=1 WHERE id=$1 AND status_id!=1", orderId)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetProductsByOrderId(orderId string) ([]model.ProductInOrder, error) {
	prods, err := ori.q.Query(
		"SELECT uuid, num, price_per_one FROM products_in_orders p WHERE order_id = $1", orderId,
	)
	if err != nil {
		return nil, err
	}
	defer prods.Close()

	products := make([]model.ProductInOrder, 0)

	for prods.Next() {
		var product model.ProductInOrder
		prods.Scan(&product.Uuid, &product.Num, &product.PricePerOne)
		products = append(products, product)
	}

	return products, nil
}

func (ori *OrderRepoImpl) GetPaymentsByOrderId(orderId string) ([]model.Payment, error) {
	pays, err := ori.q.Query(
		"SELECT total, change FROM payments p WHERE order_id = $1", orderId,
	)
	if err != nil {
		return nil, err
	}
	defer pays.Close()

	payments := make([]model.Payment, 0)

	for pays.Next() {
		var payment model.Payment
		pays.Scan(&payment.Total, &payment.Change)
		payments = append(payments, payment)
	}

	return payments, nil
}

func (ori *OrderRepoImpl) DeleteProduct(p *model.ProductInOrder) error {

	_, err := ori.q.Query(
		"DELETE FROM products_in_orders WHERE uuid = $1 AND price_per_one = $2 AND order_id = $3",
		p.Uuid, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetProductId(p *model.ProductInOrder) (int, error) {
	prodId, err := ori.q.Query(
		"SELECT id FROM products_in_orders WHERE uuid = $1 AND price_per_one = $2 AND order_id = $3",
		p.Uuid, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		return 0, err
	}
	defer prodId.Close()

	if prodId.Next() {
		var productId int
		prodId.Scan(&productId)
		return productId, nil
	}

	return -1, nil

}

func (ori *OrderRepoImpl) GetOrderStatus(orderId string) (uint8, error) {
	status, err := ori.q.Query(
		"SELECT status_id FROM orders WHERE order_id=$1",
		orderId,
	)
	if err != nil {
		return 0, err
	}
	defer status.Close()

	if status.Next() {
		var statusId uint8
		status.Scan(&statusId)
		return statusId, nil
	}

	return 0, nil

}

func (ori *OrderRepoImpl) UpdateProductNumById(num uint, id uint) (*model.ProductInOrder, error) {
	res, err := ori.q.Query(
		"UPDATE products_in_orders SET num = num + $1 WHERE id = $2"+
			"RETURNING uuid, num, price_per_one",
		num, id,
	)
	if err != nil {
		return nil, err
	}
	if res.Next() {
		var updatedProduct model.ProductInOrder
		res.Scan(&updatedProduct.Uuid, &updatedProduct.Num, &updatedProduct.PricePerOne)
		return &updatedProduct, nil
	}
	return nil, nil
}

func (ori *OrderRepoImpl) AddProduct(p *model.ProductInOrder) error {
	_, err := ori.q.Query(
		"INSERT INTO products_in_orders(uuid, num, price_per_one, order_id) VALUES($1, $2, round($3, 4), $4)",
		p.Uuid, p.Num, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) AddPayment(p *model.Payment) error {
	_, err := ori.q.Query(
		"INSERT INTO payments(total, change, order_id) VALUES($1, $2, $3)",
		p.Total, p.Change, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetById(orderId string) (*model.Order, error) {
	res, err := ori.q.Query(
		"SELECT o.id, o.short, o.date, s.name FROM orders o "+
			"LEFT JOIN statuses s ON o.status_id=s.id "+
			"WHERE o.id = $1",
		orderId,
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if !res.Next() {
		return nil, nil
	}

	var order model.Order
	res.Scan(&order.Id, &order.Short, &order.Date, &order.Status)
	order.Products, err = ori.GetProductsByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	order.Payments, err = ori.GetPaymentsByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (ori *OrderRepoImpl) GetAll(from time.Time, to time.Time) ([]model.Order, error) {
	res, err := ori.q.Query(
		"SELECT o.id, o.date, o.short, s.name FROM orders o "+
			"LEFT JOIN statuses s ON o.status_id=s.id "+
			"WHERE o.date BETWEEN $1 AND $2",
		from.Format("2006-01-02"), to.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	orders := make([]model.Order, 0)
	for res.Next() {
		var order model.Order
		res.Scan(&order.Id, &order.Date, &order.Short, &order.Status)
		order.Payments, err = ori.GetPaymentsByOrderId(order.Id)
		if err != nil {
			return nil, err
		}
		order.Products, err = ori.GetProductsByOrderId(order.Id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (ori *OrderRepoImpl) Create() (string, uint, error) {
	res, err := ori.q.Query("INSERT INTO orders(date, status_id) VALUES ($1,$2) RETURNING id, short", time.Now(), model.CREATED)
	if err != nil {
		return "", 0, err
	}
	defer res.Close()
	var id string
	var short uint
	if res.Next() {
		res.Scan(&id, &short)
	}
	return id, short, nil
}

func (ori *OrderRepoImpl) Begin() error {
	tx, err := ori.db.Begin()
	ori.q = tx
	return err
}

func (ori *OrderRepoImpl) Rollback() {
	ori.q.(*sql.Tx).Rollback()
	ori.q = ori.db
}

func (ori *OrderRepoImpl) Commit() error {
	err := ori.q.(*sql.Tx).Commit()
	ori.q = ori.db
	return err
}
