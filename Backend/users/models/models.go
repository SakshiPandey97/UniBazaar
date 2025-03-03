package models

import "gorm.io/gorm"

type Models struct {
	UserModel UserModel
	//SessionModel SessionModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		UserModel: UserModel{DB: db},
		//SessionModel: SessionModel{db: db},
	}
}
