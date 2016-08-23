package main

func (stc *Stc) ComputersinFeature(id int) ([]*Computer, error) {
	result := []*Computer{}

	rows, err := stc.Db.Query("SELECT "+
		"computer from appearance WHERE feature=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var compid int
		err = rows.Scan(&compid)
		if err != nil {
			return nil, err
		}

		c, err := stc.LoadComputer(compid)
		if err != nil {
			return nil, err
		}

		result = append(result, c)
	}

	return result, nil
}
