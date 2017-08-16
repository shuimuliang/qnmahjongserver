package dao

import "time"

func AgIDsFromAgAccount(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ag_id FROM mj.ag_account`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		agID := int32(0)

		// scan
		err = q.Scan(&agID)
		if err != nil {
			return nil, err
		}

		res = append(res, agID)
	}

	return res, nil
}

func AgUpperIDsFromAgAuth(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(ag_upper_id) FROM mj.ag_auth`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		agUpperID := int32(0)

		// scan
		err = q.Scan(&agUpperID)
		if err != nil {
			return nil, err
		}

		res = append(res, agUpperID)
	}

	return res, nil
}

func StartTimesFromAgBill(db XODB) ([]time.Time, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(start_time) FROM mj.ag_bill`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []time.Time{}
	for q.Next() {
		var startTime time.Time

		// scan
		err = q.Scan(&startTime)
		if err != nil {
			return nil, err
		}

		res = append(res, startTime)
	}

	return res, nil
}

func EmailsFromAccount(db XODB) ([]string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT email FROM mj.account`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []string{}
	for q.Next() {
		email := ""

		// scan
		err = q.Scan(&email)
		if err != nil {
			return nil, err
		}

		res = append(res, email)
	}

	return res, nil
}

func MjTypesFromCost(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(mj_type) FROM mj.cost`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		mjType := int32(0)

		// scan
		err = q.Scan(&mjType)
		if err != nil {
			return nil, err
		}

		res = append(res, mjType)
	}

	return res, nil
}

func ChannelsFromGame(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(channel) FROM mj.game`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		channel := int32(0)

		// scan
		err = q.Scan(&channel)
		if err != nil {
			return nil, err
		}

		res = append(res, channel)
	}

	return res, nil
}

func ModulesFromModule(db XODB) ([]string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT module FROM mj.module`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []string{}
	for q.Next() {
		module := ""

		// scan
		err = q.Scan(&module)
		if err != nil {
			return nil, err
		}

		res = append(res, module)
	}

	return res, nil
}

func PmsnTypesFromPermission(db XODB) ([]string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(pmsn_type) FROM mj.permission`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []string{}
	for q.Next() {
		pmsnType := ""

		// scan
		err = q.Scan(&pmsnType)
		if err != nil {
			return nil, err
		}

		res = append(res, pmsnType)
	}

	return res, nil
}

func RolesFromRole(db XODB) ([]string, error) {
	var err error

	// sql query
	const sqlstr = `SELECT role FROM mj.role`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []string{}
	for q.Next() {
		role := ""

		// scan
		err = q.Scan(&role)
		if err != nil {
			return nil, err
		}

		res = append(res, role)
	}

	return res, nil
}

func ChannelsFromShop(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT distinct(channel) FROM mj.shop`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		channel := int32(0)

		// scan
		err = q.Scan(&channel)
		if err != nil {
			return nil, err
		}

		res = append(res, channel)
	}

	return res, nil
}

func PlayerIDsFromPlayer(db XODB) ([]int32, error) {
	var err error

	// sql query
	const sqlstr = `SELECT player_id FROM mj.player`

	// run query
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}

	// load results
	res := []int32{}
	for q.Next() {
		playerID := int32(0)

		// scan
		err = q.Scan(&playerID)
		if err != nil {
			return nil, err
		}

		res = append(res, playerID)
	}

	return res, nil
}
