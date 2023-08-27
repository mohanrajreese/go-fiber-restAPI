package models

type Books struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Language  *string `json:"language"`
	Publisher *string `json:"publisher"`
	Price     *int    `json:"price"`
}
