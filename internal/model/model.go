package model

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}
