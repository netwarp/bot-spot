package database

import (
	"database/sql"
	"fmt"
)

type Status string

const (
	Buy       Status = "buy"
	Sell      Status = "sell"
	Completed Status = "completed"
)

type BuyStruct struct {
	Offset int
	Price  float64
	ID     string
}

type SellStruct struct {
	Offset int
	Price  float64
	ID     string
}

type MetaData struct {
	FreeBalanceUSD float64
	USDDedicated   float64
	Percent        int
	BTCPrice       float64
}

type Cycle struct {
	Id       int
	Exchange string
	Status   Status
	Quantity float64
	Buy      BuyStruct
	Sell     SellStruct
	MetaData MetaData
}

func CycleNew(cycle *Cycle) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return 0, fmt.Errorf("error getting database: %v", err)
	}

	id, err := db.Exec("INSERT INTO cycles (exchange, status, quantity, buyPrice, buyId, sellPrice, sellId) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id", cycle.Exchange, cycle.Status, cycle.Quantity, cycle.Buy.Price, cycle.Buy.ID, cycle.Sell.Price, cycle.Sell.ID)
	if err != nil {
		return 0, fmt.Errorf("error inserting cycle: %v", err)
	}

	insertId, err := id.LastInsertId()
	if err != nil {
		return 0, err
	}

	defer db.Close()

	return insertId, nil
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
