package database

import (
	"database/sql"
)

type Status string

const (
	Buy       Status = "buy"
	Sell      Status = "sell"
	Completed Status = "completed"
)

type BuyStruct struct {
	Price float64
	ID    string
}

type SellStruct struct {
	Price float64
	ID    string
}

type Cycle struct {
	Id       int
	Exchange string
	Status   Status
	Quantity float64
	Buy      BuyStruct
	Sell     SellStruct
}

func CycleNew(cycle *Cycle) (sql.Result, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec("INSERT INTO cycles (exchange, status, quantity, buyPrice, buyId, sellPrice, sellId) VALUES (?, ?, ?, ?, ?, ?, ?)", cycle.Exchange, cycle.Status, cycle.Quantity, cycle.Buy.Price, cycle.Buy.ID, cycle.Sell.Price, cycle.Sell.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CycleList() ([]Cycle, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	var cycles []Cycle

	rows, err := db.Query("SELECT * FROM cycles ORDER BY id DESC ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cycle Cycle
		err := rows.Scan(&cycle.Id, &cycle.Exchange, &cycle.Status, &cycle.Quantity, &cycle.Buy.Price, &cycle.Buy.ID, &cycle.Sell.Price, &cycle.Sell.ID)
		if err != nil {
			return nil, err
		}
		cycles = append(cycles, cycle)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cycles, nil
}

func CycleGetById(id int) (*Cycle, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT * FROM cycles WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var cycle Cycle
	err = rows.Scan(&cycle.Id, &cycle.Exchange, &cycle.Status, &cycle.Quantity, &cycle.Buy.Price, &cycle.Buy.ID, &cycle.Sell.Price, &cycle.Sell.ID)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &cycle, nil
}

func CycleDeleteById(id int) error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM cycles WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func CycleListPerPage(page int, itemsPerPage int) ([]Cycle, error) {

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	skip := (page - 1) * itemsPerPage
	rows, err := db.Query("SELECT * FROM cycles ORDER BY id DESC LIMIT ? OFFSET ?", itemsPerPage, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cycles []Cycle
	for rows.Next() {
		var cycle Cycle
		err := rows.Scan(&cycle.Id, &cycle.Exchange, &cycle.Status, &cycle.Quantity, &cycle.Buy.Price, &cycle.Buy.ID, &cycle.Sell.Price, &cycle.Sell.ID)
		if err != nil {
			return nil, err
		}
		cycles = append(cycles, cycle)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cycles, nil
}

func CycleUpdate(id int, field string, value interface{}) (sql.Result, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec("UPDATE cycles SET "+field+" = ? WHERE id = ?", value, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
