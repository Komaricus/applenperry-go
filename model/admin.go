package model

type AdminsTable struct{}

func (AdminsTable) TableName() string {
	return "dbo.admins"
}

type Admin struct {
	AdminsTable
	ID       string `json:"id" gorm:"primarykey"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
