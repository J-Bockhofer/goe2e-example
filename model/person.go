package model

type Person struct {
	Name string `json:"name" binding:"required,min=2"`
	Age  int    `json:"age" binding:"required"`
}
