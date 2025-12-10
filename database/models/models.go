package models

type Student struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Gpa float64 `json:"gpa"`
	Active bool `json:"active"`
}

type Record struct {
	Id int
	Student Student
}

type RecordInfo struct {
	Name   string
	Gpa    float64
	Active bool
}