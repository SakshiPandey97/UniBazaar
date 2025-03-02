package models

import "gorm.io/gorm"

type Models struct {
	UserModel UserModel
	//SessionModel SessionModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		UserModel: UserModel{db: db},
		//SessionModel: SessionModel{db: db},
	}
}
