package main

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var db sql.DB

func dbConnect(dbConnStr string) sql.DB {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		fmt.Println(err)
	}

	return *db
}

func dbWriteOrder(m Model) {

	deliveryID := dbWriteDelivery(m)
	paymentId := dbWritePayment(m)
	itemsIds := dbWriteItems(m)
	var query = fmt.Sprintf("INSERT INTO orders VALUES ('%v','%v','%v',%v,'%v',$1,'%v','%v','%v','%v','%v',%v,$2,'%v')",
		m.OrderUid,
		m.TrackNumber,
		m.Entry,
		deliveryID,
		paymentId,
		m.Locale,
		m.InternalSignature,
		m.CustomerId,
		m.DeliveryService,
		m.Shardkey,
		m.SmId,
		m.OofShard)

	_, err := db.Exec(query, pq.Array(itemsIds), pq.FormatTimestamp(m.DateCreated))
	if err != nil {
		fmt.Println(err.Error())
	}

}

func dbWriteDelivery(m Model) int {

	var id int
	var query = fmt.Sprintf("INSERT INTO deliveries VALUES (default, '%v','%v','%v','%v','%v','%v','%v') RETURNING id",
		m.Delivery.Name,
		m.Delivery.Phone,
		m.Delivery.Zip,
		m.Delivery.City,
		m.Delivery.Address,
		m.Delivery.Region,
		m.Delivery.Email)

	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return id
}

func dbWritePayment(m Model) string {

	checkQuery := fmt.Sprintf("SELECT 1 FROM payments WHERE transaction='%v'", m.Payment.Transaction)
	rows, _ := db.Query(checkQuery)
	if !rows.Next() {
		var transaction string
		var query = fmt.Sprintf("INSERT INTO payments VALUES ('%v','%v','%v','%v',%v,%v,'%v',%v,%v,%v) RETURNING transaction",
			m.Payment.Transaction,
			m.Payment.RequestId,
			m.Payment.Currency,
			m.Payment.Provider,
			m.Payment.Amount,
			m.Payment.PaymentDt,
			m.Payment.Bank,
			m.Payment.DeliveryCost,
			m.Payment.GoodsTotal,
			m.Payment.CustomFee)

		err := db.QueryRow(query).Scan(&transaction)
		if err != nil {
			fmt.Println(err.Error())
		}

		return transaction
	}
	return m.Payment.Transaction

}

func dbWriteItems(m Model) []int {
	var ids []int
	for _, item := range m.Items {
		var id int
		var query = fmt.Sprintf("INSERT INTO items VALUES (default, '%v',%v,'%v','%v',%v,'%v',%v,%v,'%v',%v) RETURNING chrt_id",
			item.TrackNumber,
			item.Price, item.Rid,
			item.Name, item.Sale,
			item.Size,
			item.Total_price,
			item.Nm_id,
			item.Brand,
			item.Status)

		err := db.QueryRow(query).Scan(&id)
		if err != nil {
			fmt.Println(err.Error())
		}
		ids = append(ids, id)
	}
	return ids
}

func dbGetOrder(uid string) Model {

	var result Model

	var (
		deliveryId int
		paymentId  string
		itemsIds   pq.Int64Array
	)

	query := fmt.Sprintf("SELECT * FROM orders WHERE order_uid='%v'", uid)

	err := db.QueryRow(query).Scan(&result.OrderUid,
		&result.TrackNumber,
		&result.Entry,
		&deliveryId,
		&paymentId,
		&itemsIds,
		&result.Locale,
		&result.InternalSignature,
		&result.CustomerId,
		&result.DeliveryService,
		&result.Shardkey,
		&result.SmId,
		&result.DateCreated,
		&result.OofShard)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get delivery by ID
	query = fmt.Sprintf("SELECT * FROM deliveries WHERE id=%v", deliveryId)

	err = db.QueryRow(query).Scan(&deliveryId,
		&result.Delivery.Name,
		&result.Delivery.Phone,
		&result.Delivery.Zip,
		&result.Delivery.City,
		&result.Delivery.Address,
		&result.Delivery.Region,
		&result.Delivery.Email)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get payment by transaction
	query = fmt.Sprintf("SELECT * FROM payments WHERE transaction='%v'", paymentId)
	err = db.QueryRow(query).Scan(&result.Payment.Transaction,
		&result.Payment.RequestId,
		&result.Payment.Currency,
		&result.Payment.Provider,
		&result.Payment.Amount,
		&result.Payment.PaymentDt,
		&result.Payment.Bank,
		&result.Payment.DeliveryCost,
		&result.Payment.GoodsTotal,
		&result.Payment.CustomFee)
	if err != nil {
		fmt.Println(err.Error())
	}

	//get each item by []ids
	for _, i := range itemsIds {
		var item Item
		query = fmt.Sprintf("SELECT * FROM items WHERE chrt_id=%v", i)
		err = db.QueryRow(query).Scan(&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.Total_price,
			&item.Nm_id,
			&item.Brand,
			&item.Status)
		if err != nil {
			fmt.Println(err.Error())
		}

		result.Items = append(result.Items, item)
	}

	return result
}
