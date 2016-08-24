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

func (stc *Stc) FeaturesForComputer(id int) ([]*Feature, error) {
	result := []*Feature{}

	rows, err := stc.Db.Query("SELECT "+
		"feature from appearance WHERE computer=?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var featid int
		err = rows.Scan(&featid)
		if err != nil {
			return nil, err
		}

		f, err := stc.LoadFeature(featid)
		if err != nil {
			return nil, err
		}

		result = append(result, f)
	}

	return result, nil
}
