package example

import "time"

type Example struct {
	Id         int       `json:"id" gorm:"column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	Age        int       `json:"age" gorm:"column:age"`
	Gender     string    `json:"gender" gorm:"column:gender"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}
