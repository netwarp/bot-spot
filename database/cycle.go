package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
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
	Percent        float64
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
	defer func() { _ = db.Close() }()

	// Retry INSERT on transient SQLITE_BUSY/database is locked errors
	var res sql.Result
	for attempt := 0; attempt < 5; attempt++ {
		res, err = db.Exec("INSERT INTO cycles (exchange, status, quantity, buyPrice, buyId, sellPrice, sellId, freeBalance, dedicatedBalance, buyOffset, sellOffset, percent, btcPrice) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id", cycle.Exchange, cycle.Status, cycle.Quantity, cycle.Buy.Price, cycle.Buy.ID, cycle.Sell.Price, cycle.Sell.ID, cycle.MetaData.FreeBalanceUSD, cycle.MetaData.USDDedicated, cycle.Buy.Offset, cycle.Sell.Offset, cycle.MetaData.Percent, cycle.MetaData.BTCPrice)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), "SQLITE_BUSY") || strings.Contains(err.Error(), "database is locked") {
			time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond)
			continue
		}
		return 0, fmt.Errorf("error inserting cycle: %v", err)
	}
	if err != nil {
		return 0, fmt.Errorf("error inserting cycle after retries: %v", err)
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertId, nil
}

func CycleList() ([]Cycle, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	var cycles []Cycle

	rows, err := db.Query("SELECT * FROM cycles ORDER BY id DESC ")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cycle Cycle
		err := rows.Scan(
			&cycle.Id,
			&cycle.Exchange,
			&cycle.Status,
			&cycle.Quantity,
			&cycle.Buy.Price,
			&cycle.Buy.ID,
			&cycle.Sell.Price,
			&cycle.Sell.ID,
			&cycle.MetaData.FreeBalanceUSD,
			&cycle.MetaData.USDDedicated,
			&cycle.Buy.Offset,
			&cycle.Sell.Offset,
			&cycle.MetaData.Percent,
			&cycle.MetaData.BTCPrice,
		)
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
	defer func() { _ = db.Close() }()

	rows, err := db.Query("SELECT * FROM cycles WHERE id = ? LIMIT 1", id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if !rows.Next() {
		if err = rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}

	var cycle Cycle
	err = rows.Scan(&cycle.Id, &cycle.Exchange, &cycle.Status, &cycle.Quantity, &cycle.Buy.Price, &cycle.Buy.ID, &cycle.Sell.Price, &cycle.Sell.ID, &cycle.MetaData.FreeBalanceUSD, &cycle.MetaData.USDDedicated, &cycle.Buy.Offset, &cycle.Sell.Offset, &cycle.MetaData.Percent, &cycle.MetaData.BTCPrice)
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
	defer func() { _ = db.Close() }()

	for attempt := 0; attempt < 5; attempt++ {
		_, err = db.Exec("DELETE FROM cycles WHERE id = ?", id)
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "SQLITE_BUSY") || strings.Contains(err.Error(), "database is locked") {
			time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond)
			continue
		}
		return err
	}
	return err
}

func CycleListPerPage(page int, itemsPerPage int) ([]Cycle, error) {

	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	defer func() { _ = db.Close() }()

	skip := (page - 1) * itemsPerPage
	rows, err := db.Query("SELECT * FROM cycles ORDER BY id DESC LIMIT ? OFFSET ?", itemsPerPage, skip)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var cycles []Cycle
	for rows.Next() {
		var cycle Cycle
		err = rows.Scan(&cycle.Id, &cycle.Exchange, &cycle.Status, &cycle.Quantity, &cycle.Buy.Price, &cycle.Buy.ID, &cycle.Sell.Price, &cycle.Sell.ID, &cycle.MetaData.FreeBalanceUSD, &cycle.MetaData.USDDedicated, &cycle.Buy.Offset, &cycle.Sell.Offset, &cycle.MetaData.Percent, &cycle.MetaData.BTCPrice)
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
	defer func() {
		cerr := db.Close()
		if cerr != nil {
			log.Printf("warning: closing db in CycleUpdate: %v", cerr)
		}
	}()

	var result sql.Result
	for attempt := 0; attempt < 5; attempt++ {
		result, err = db.Exec("UPDATE cycles SET "+field+" = ? WHERE id = ?", value, id)
		if err == nil {
			return result, nil
		}
		if strings.Contains(err.Error(), "SQLITE_BUSY") || strings.Contains(err.Error(), "database is locked") {
			time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond)
			continue
		}
		return nil, err
	}
	return nil, err
}

// helpers
func (c *Cycle) CalcPercent() float64 {
	totalBuy := c.Buy.Price * c.Quantity
	totalSell := c.Sell.Price * c.Quantity

	percent := (totalSell - totalBuy) / totalBuy * 100
	return percent
}

func (c *Cycle) CalcProfit() float64 {
	totalBuy := c.Buy.Price * c.Quantity
	totalSell := c.Sell.Price * c.Quantity

	profit := totalSell - totalBuy
	return profit
}

// String returns a detailed string representation of a Cycle, useful for logs.
func (c Cycle) String() string {
	return fmt.Sprintf(
		"Cycle{id:%d, ex:%s, status:%s, qty:%.8f, buy:{off:%d price:%.8f id:%s}, sell:{off:%d price:%.8f id:%s}, meta:{freeUSD:%.2f dedicatedUSD:%.2f percent:%d btc:%.2f}, profit:%.8f, pct:%.4f%%}",
		c.Id,
		c.Exchange,
		c.Status,
		c.Quantity,
		c.Buy.Offset,
		c.Buy.Price,
		c.Buy.ID,
		c.Sell.Offset,
		c.Sell.Price,
		c.Sell.ID,
		c.MetaData.FreeBalanceUSD,
		c.MetaData.USDDedicated,
		c.MetaData.Percent,
		c.MetaData.BTCPrice,
		c.CalcProfit(),
		c.CalcPercent(),
	)
}
