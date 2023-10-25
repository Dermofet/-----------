package view

type UserView struct {
	ID    string `json:"id"`    // ID
	Name  string `json:"name"`  // Имя в формате ФИО
	Age   int    `json:"age"`   // Возраст
	Email string `json:"email"` // Электронная почта
	Phone string `json:"phone"` // Номер мобильного телефона
}
