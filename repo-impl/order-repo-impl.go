package repoimpl

import (
	"database/sql"
	"order-and-pay/model"
	"order-and-pay/utils"
	"time"
)

type OrderRepoImpl struct {
	db *sql.DB
}

func NewOrderRepoImpl(db *sql.DB) *OrderRepoImpl {
	return &OrderRepoImpl{db}
}

func (ori *OrderRepoImpl) GetPaymentsSumByOrderId(orderId string) float32 {
	payRes, err := ori.db.Query("SELECT SUM(total)-SUM(change) FROM payments WHERE order_id=$1", orderId)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer payRes.Close()
	var paymentSum float32
	if payRes.Next() {
		payRes.Scan(&paymentSum)
	}
	return paymentSum
}

func (ori *OrderRepoImpl) GetProductsPriceSumByOrderId(orderId string) float32 {
	prodPriceSumRes, err := ori.db.Query("SELECT SUM(num*price_per_one) FROM products_in_orders WHERE order_id=$1", orderId)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer prodPriceSumRes.Close()
	var priceSum float32
	if prodPriceSumRes.Next() {
		prodPriceSumRes.Scan(&priceSum)
	}
	return priceSum
}

func (ori *OrderRepoImpl) UpdateOrderStatusToComplete(orderId string) {
	_, err := ori.db.Query("UPDATE orders SET status_id=1 WHERE id=$1", orderId)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
}

func (ori *OrderRepoImpl) GetProductsByOrderId(orderId string) []model.ProductInOrder {
	prods, err := ori.db.Query(
		"SELECT uuid, num, price_per_one FROM products_in_orders p WHERE order_id = $1", orderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer prods.Close()

	products := make([]model.ProductInOrder, 0)

	for prods.Next() {
		var product model.ProductInOrder
		prods.Scan(&product.Uuid, &product.Num, &product.PricePerOne)
		products = append(products, product)
	}

	return products
}

func (ori *OrderRepoImpl) GetPaymentsByOrderId(orderId string) []model.Payment {
	pays, err := ori.db.Query(
		"SELECT total, change FROM payments p WHERE order_id = $1", orderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer pays.Close()

	payments := make([]model.Payment, 0)

	for pays.Next() {
		var payment model.Payment
		pays.Scan(&payment.Total, &payment.Change)
		payments = append(payments, payment)
	}

	return payments
}

func (ori *OrderRepoImpl) DeleteProduct(p *model.ProductInOrder) {
	_, err := ori.db.Query(
		"DELETE FROM products_in_orders WHERE uuid = $1 AND price_per_one = $2 AND order_id = $3",
		p.Uuid, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
}

func (ori *OrderRepoImpl) CheckAndGetProductId(p *model.ProductInOrder) (uint, bool) {
	prodId, err := ori.db.Query(
		"SELECT id FROM products_in_orders WHERE uuid = $1 AND price_per_one = $2 AND order_id = $3",
		p.Uuid, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer prodId.Close()

	if prodId.Next() {
		var productId uint
		prodId.Scan(&productId)
		return productId, true
	}

	return 0, false

}

func (ori *OrderRepoImpl) UpdateProductNumById(num uint, id uint) {
	_, err := ori.db.Query("UPDATE products_in_orders SET num = num + $1 WHERE id = $2", num, id)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
}

func (ori *OrderRepoImpl) AddProduct(p *model.ProductInOrder) {
	_, err := ori.db.Query(
		"INSERT INTO products_in_orders(uuid, num, price_per_one, order_id) VALUES($1, $2, round($3, 4), $4)",
		p.Uuid, p.Num, p.PricePerOne, p.OrderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
}

func (ori *OrderRepoImpl) AddPayment(p *model.Payment) {
	_, err := ori.db.Query("INSERT INTO payments(total, change, order_id) VALUES($1, $2, $3)", p.Total, p.Change, p.OrderId)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
}

func (ori *OrderRepoImpl) GetById(orderId string) *model.Order {
	res, err := ori.db.Query(
		"SELECT o.id, o.short, o.date, s.name FROM orders o "+
			"LEFT JOIN statuses s ON o.status_id=s.id "+
			"WHERE o.id = $1",
		orderId,
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer res.Close()

	if !res.Next() {
		return nil
	}

	var order model.Order
	res.Scan(&order.Id, &order.Short, &order.Date, &order.Status)
	order.Products = ori.GetProductsByOrderId(orderId)
	order.Payments = ori.GetPaymentsByOrderId(orderId)
	return &order
}

func (ori *OrderRepoImpl) GetAll(from time.Time, to time.Time) []model.Order {
	res, err := ori.db.Query(
		"SELECT o.id, o.date, o.short, s.name FROM orders o "+
			"LEFT JOIN statuses s ON o.status_id=s.id "+
			"WHERE o.date BETWEEN $1 AND $2",
		from.Format("2006-01-02"), to.Format("2006-01-02"),
	)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer res.Close()

	orders := make([]model.Order, 0)
	for res.Next() {
		var order model.Order
		res.Scan(&order.Id, &order.Date, &order.Short, &order.Status)
		order.Payments = ori.GetPaymentsByOrderId(order.Id)
		order.Products = ori.GetProductsByOrderId(order.Id)
		orders = append(orders, order)
	}
	return orders
}

func (ori *OrderRepoImpl) Create() (string, uint) {
	res, err := ori.db.Query("INSERT INTO orders(date, status_id) VALUES ($1,$2) RETURNING id, short", time.Now(), model.CREATED)
	if err != nil {
		utils.Logger.Panic(err.Error())
	}
	defer res.Close()
	var id string
	var short uint
	if res.Next() {
		res.Scan(&id, &short)
	}
	return id, short
}
