package dao

func (g *Game) SetExist(exist bool) {
	g._exists = exist
}

func (c *Cost) SetExist(exist bool) {
	c._exists = exist
}

func (m *Module) SetExist(exist bool) {
	m._exists = exist
}

func (s *Shop) SetExist(exist bool) {
	s._exists = exist
}

func (a *Account) SetExist(exist bool) {
	a._exists = exist
}

func (r *Role) SetExist(exist bool) {
	r._exists = exist
}

func (p *Permission) SetExist(exist bool) {
	p._exists = exist
}

func SelectAllGames(db XODB) ([]*Game, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.game`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Game{}
	for q.Next() {
		g := Game{
			_exists: true,
		}

		// scan
		err = q.Scan(&g.IndexID, &g.Channel, &g.Version, &g.Size, &g.Module, &g.MjTypes, &g.Enabled, &g.UpdateType, &g.DownloadURL, &g.SvnVersion)
		if err != nil {
			return nil, err
		}

		res = append(res, &g)
	}

	return res, nil
}

func SelectAllCosts(db XODB) ([]*Cost, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.cost`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Cost{}
	for q.Next() {
		c := Cost{
			_exists: true,
		}

		// scan
		err = q.Scan(&c.IndexID, &c.MjType, &c.MjDesc, &c.Rounds, &c.Cards, &c.Coins)
		if err != nil {
			return nil, err
		}

		res = append(res, &c)
	}

	return res, nil
}

func SelectAllModules(db XODB) ([]*Module, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.module`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Module{}
	for q.Next() {
		m := Module{
			_exists: true,
		}

		// scan
		err = q.Scan(&m.IndexID, &m.Module, &m.Comment)
		if err != nil {
			return nil, err
		}

		res = append(res, &m)
	}

	return res, nil
}

func SelectAllShops(db XODB) ([]*Shop, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.shop`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Shop{}
	for q.Next() {
		s := Shop{
			_exists: true,
		}

		// scan
		err = q.Scan(&s.IndexID, &s.Channel, &s.PayType, &s.GemID, &s.WaresID, &s.WaresName, &s.GoodsCount, &s.ExtraCount, &s.Price, &s.IconURL)
		if err != nil {
			return nil, err
		}

		res = append(res, &s)
	}

	return res, nil
}

func SelectAllAccounts(db XODB) ([]*Account, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.account`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Account{}
	for q.Next() {
		a := Account{
			_exists: true,
		}

		// scan
		err = q.Scan(&a.IndexID, &a.Email, &a.Password, &a.Role)
		if err != nil {
			return nil, err
		}

		res = append(res, &a)
	}

	return res, nil
}

func SelectAllRoles(db XODB) ([]*Role, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.role`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Role{}
	for q.Next() {
		r := Role{
			_exists: true,
		}

		// scan
		err = q.Scan(&r.IndexID, &r.Role, &r.Comment, &r.Permissions)
		if err != nil {
			return nil, err
		}

		res = append(res, &r)
	}

	return res, nil
}

func SelectAllPermissions(db XODB) ([]*Permission, error) {
	var err error

	// sql query
	const sqlstr = `SELECT * FROM mj.permission`

	// run query
	XOLog(sqlstr)
	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*Permission{}
	for q.Next() {
		p := Permission{
			_exists: true,
		}

		// scan
		err = q.Scan(&p.IndexID, &p.PmsnType, &p.PmsnContent, &p.Comment)
		if err != nil {
			return nil, err
		}

		res = append(res, &p)
	}

	return res, nil
}
