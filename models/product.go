package models

func GetAllRobot() ([]*Robot, error) {
	var robots []*Robot
	err := db.Find(&robots).Error
	if err != nil {
		return nil, err
	}

	return robots, nil
}
