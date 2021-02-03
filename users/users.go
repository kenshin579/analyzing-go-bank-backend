package users

import (
	"fmt"
	"time"

	"duomly.com/go-bank-backend/helpers"
	"duomly.com/go-bank-backend/interfaces"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Login(username string, pass string) map[string]interface{} {
	// Connect DB
	db := helpers.ConnectDB()
	user := &interfaces.User{}

	//SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((username = 'Martin' )) ORDER BY "users"."id" ASC LIMIT 1
	db.LogMode(true)
	if db.Where("username = ? ", username).First(&user).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"} //db 조회시 없는 경우
	}
	// Verify password
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return map[string]interface{}{"message": "Wrong password"}
	}
	// Find accounts for the user
	accounts := []interfaces.ResponseAccount{}

	//SELECT id, name, balance FROM "accounts"  WHERE (user_id = 1 )
	db.Table("accounts").Select("id, name, balance").Where("user_id = ? ", user.ID).Scan(&accounts)

	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	defer db.Close()

	//token에는 어떤 내용을 담는가? (userId, IssueAt, expiryDate, 사용한 알고리즘)
	//
	// Sign token
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	// Prepare response
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	fmt.Printf("response : %v\n", response)
	return response
}
