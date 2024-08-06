package user

import (
	"net/http"
	"optimaHurt/hurtownie"
)

type User struct {
	Client *http.Client
	Hurts  []hurtownie.IHurt
	Creds  []UserCreds
}

func (u *User) TakeHurtCreds(name hurtownie.HurtName) UserCreds {
	for _, i := range u.Creds {
		if i.HurtName == name {
			return i
		}
	}
	return UserCreds{}
}

type DataBaseUserObject struct {
	Email          string      `bson:"email" json:"email"`
	Username       string      `bson:"username" json:"username"`
	Password       string      `bson:"password" json:"password"`
	CompanyData    CompanyData `bson:"companyData" json:"companyData"`
	AvailableHurts int         `bson:"availableHurts" json:"availableHurts"`
	Creds          []UserCreds `bson:"creds" json:"creds"`
}

type UserCreds struct {
	HurtName hurtownie.HurtName `json:"hurtName" bson:"hurtName"`
	Login    string             `json:"login" bson:"login"`
	Password string             `json:"password" bson:"password"`
}

type Adress struct {
	Street string `bson:"street"  bson:"street"`
	Nr     int    `bson:"nr" json:"nr"`
}

type CompanyData struct {
	Name   string `bson:"CompanyName" json:"companyName"`
	Nip    string `bson:"nip" json:"nip"`
	Adress Adress `bson:"adress" json:"adress"`
}
