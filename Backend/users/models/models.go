package models

import "gorm.io/gorm"

type Models struct {
	UserModel UserModel
}

func NewModels(conn *gorm.DB) Models {
	return Models{
		UserModel: UserModel{db: conn},
	}

}
