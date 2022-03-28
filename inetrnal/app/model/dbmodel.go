package model

type DbModelRequest struct {
	Parameter1 string `json:"parameter1" validate:"required"`
	Parameter2 int    `json:"parameter2" validate:"required,numeric,min=1"`
}

type DbModel struct {
	FirstField  string `db:"first_field" json:"firstField"`
	SecondField int    `db:"second_field" json:"secondField"`
}
