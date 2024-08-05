package account

import (
	"github.com/gin-gonic/gin"
	. "optimaHurt/constAndVars"
	"optimaHurt/hurtownie"
	. "optimaHurt/user"
)

func AddUser(c *gin.Context) {
	con := DbConnect.Collection(UserCollection)
	newUser := DataBaseUserObject{
		Email:          "pog@gmail.com",
		Username:       "kholin",
		Password:       "nice",
		AvailableHurts: int(hurtownie.Eurocash + hurtownie.Sot + hurtownie.Specjal + hurtownie.Tedi),
		Creds: []UserCreds{
			{HurtName: hurtownie.Eurocash, Login: "Sp.gaj", Password: "0103Gaj1"},
			{HurtName: hurtownie.Sot, Login: "Sot.22803", Password: "perelka"},
			{HurtName: hurtownie.Specjal, Login: "21B879.GAJWILKSZ", Password: "YN38544P"},
			{HurtName: hurtownie.Tedi, Login: "lukasz@delikatesykredens.pl", Password: "dqfciavfbvuzrdsx"},
		},
	}

	res, err := con.InsertOne(ContextBackground, newUser)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"inserted": res})
}
