package dao

import (
	"fmt"
	"time"
)

func (ab *AgBill) GetDate(time time.Time) string {
	return time.Format("2006-01-02")
}

func (ab *AgBill) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

func (ap *AgPay) GetMoney(money int32) string {
	return fmt.Sprintf("%.2f", float64(money)/100)
}

func (ap *AgPay) GetTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func UpdateAgPayStatus(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE ag_pay SET delflag = '1' WHERE delflag = '0' `

	// run query
	_, err = db.Exec(sqlstr)
	return err
}

func UpdateAgBillStatus(db XODB) error {
	var err error

	// sql query
	const sqlstr = `UPDATE ag_bill SET delflag = '2' WHERE delflag = '0' and last_week_dakuan < 10000`

	// run query
	_, err = db.Exec(sqlstr)
	return err
}
