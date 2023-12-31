package repoimpl

import (
	"order-and-pay/db"
	"order-and-pay/iface"
	"order-and-pay/model"
	"time"
)

type OrderRepoImpl struct {
	db *db.SqlDb
}

func NewOrderRepoImpl(db *db.SqlDb) *OrderRepoImpl {
	return &OrderRepoImpl{db}
}

func (ori *OrderRepoImpl) GetPaymentsSumByOrderId(tx iface.Itx, orderId string) (float32, error) {
	payRes, err := tx.Query("SELECT SUM(total)-SUM(change) FROM payments WHERE order_id=$1", orderId)
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

func (ori *OrderRepoImpl) GetProductsPriceSumByOrderId(tx iface.Itx, orderId string) (float32, error) {
	prodPriceSumRes, err := tx.Query("SELECT SUM(num*price_per_one) FROM products_in_orders WHERE order_id=$1", orderId)
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

func (ori *OrderRepoImpl) UpdateOrderStatusToComplete(tx iface.Itx, orderId string) error {
	_, err := tx.Query("UPDATE orders SET status_id=2, date = NOW() WHERE id=$1 AND status_id!=2", orderId)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetProductsByOrderId(orderId string) ([]model.ProductInOrder, error) {
	prods, err := ori.db.Query(
		"SELECT uuid, num, price_per_one FROM products_in_orders WHERE order_id=$1", orderId,
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
	pays, err := ori.db.Query(
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
	_, err := ori.db.Query(
		"DELETE FROM products_in_orders WHERE uuid = $1 AND price_per_one = $2 AND order_id = $3",
		p.Uuid, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetProductId(tx iface.Itx, p *model.ProductInOrder) (int, error) {
	prodId, err := tx.Query(
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

func (ori *OrderRepoImpl) GetOrderStatus(tx iface.Itx, orderId string) (uint8, error) {
	status, err := tx.Query(
		"SELECT status_id FROM orders WHERE id=$1",
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

func (ori *OrderRepoImpl) UpdateProductNumById(tx iface.Itx, num uint, id uint) error {
	_, err := tx.Query(
		"UPDATE products_in_orders SET num = num + $1 WHERE id = $2",
		num, id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) AddProduct(tx iface.Itx, p *model.ProductInOrder) error {
	_, err := tx.Query(
		"INSERT INTO products_in_orders(uuid, num, price_per_one, order_id) VALUES($1, $2, round($3, 4), $4)",
		p.Uuid, p.Num, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) AddPayment(p *model.Payment) error {
	_, err := ori.db.Query(
		"INSERT INTO payments(total, change, order_id) VALUES($1, $2, $3)",
		p.Total, p.Change, p.OrderId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (ori *OrderRepoImpl) GetById(orderId string) (*model.Order, error) {

	order, err := ori.getRawOrderById(orderId)
	if err != nil {
		return nil, err
	}

	order.Products, err = ori.GetProductsByOrderId(orderId)
	if err != nil {
		return nil, err
	}
	order.Payments, err = ori.GetPaymentsByOrderId(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ori *OrderRepoImpl) getRawOrderById(orderId string) (*model.Order, error) {
	res, err := ori.db.Query(
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
	return &order, nil
}

func (ori *OrderRepoImpl) GetAll(from time.Time, to time.Time) ([]model.Order, error) {

	orders, err := ori.getAllRaw(from, to)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		orders[i].Payments, err = ori.GetPaymentsByOrderId(orders[i].Id)
		if err != nil {
			return nil, err
		}
		orders[i].Products, err = ori.GetProductsByOrderId(orders[i].Id)
		if err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (ori *OrderRepoImpl) getAllRaw(from time.Time, to time.Time) ([]model.Order, error) {
	res, err := ori.db.Query(
		"SELECT o.id, o.date, o.short, s.name FROM orders o "+
			"LEFT JOIN statuses s ON o.status_id=s.id "+
			"WHERE (o.date BETWEEN $1 AND $2) AND (s.name = 'COMPLETED')",
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
		orders = append(orders, order)
	}
	return orders, nil
}

func (ori *OrderRepoImpl) Create() (string, uint, error) {
	res, err := ori.db.Query("INSERT INTO orders(date, status_id) VALUES ($1,$2) RETURNING id, short", time.Now(), model.CREATED)
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

func (ori *OrderRepoImpl) GetDb() iface.Idb {
	return ori.db
}
