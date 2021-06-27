package models

func CreateAuthority(auth Authority) (id string, err error) {
	//var authorityBox Authority
	/*if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id")
	}*/
	if err := db.Create(&auth).Error; err != nil {
		return "0", err
	}

	return auth.AuthorityId, nil
}

func UpdateAuthority(auth Authority) (err error,authority Authority) {
	err = db.Where("authority_id = ?", auth.AuthorityId).First(&Authority{}).Updates(&auth).Error
	return err,auth
}

func SetAuthority(auth Authority)  error{
	var s Authority
	db.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	err := db.Model(&s).Association("AuthorityId").Replace(&auth.AuthorityId)
	return err.Error
}

func DeleteAuthority(auth Authority) (err error) {
	err = db.Unscoped().Delete(auth).Error
	return err
}