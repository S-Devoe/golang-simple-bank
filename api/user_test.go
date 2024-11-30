package api

// import (
// 	"testing"

// 	"github.com/S-Devoe/golang-simple-bank/util"
// 	"go.uber.org/mock/gomock"
// )

// func TestCreateUserAPI(t *testing.T) {
// 	user_req := createUserRequest{
// 		Username: util.GenerateRandomUserName(),
// 		Email:    util.GenerateRandomEmail(),
// 		FullName: util.GenerateRandomName(),
// 	}
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// }

// func randomAccount() db.Account {
// 	user := randomUser()
// 	return db.Account{
// 		ID:       util.GenerateRandomInt(1, 100),
// 		Owner:    user.Username,
// 		Balance:  float64(util.GenerateRandomInt(0, 100)),
// 		Currency: util.GenerateRandomCurrency(),
// 	}
// }
