package models

import (
	"fmt"
	"os"
	u "restserver/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	//gorm.Model
	ID        uint `json:"ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string `json:"email"`
	Password  string `json:"password"`
	Token     string `json:"token"`
}

//Validate incoming user details...
func (account *Account) Validate() (map[string]interface{}, bool) {
	fmt.Println("Email %", account.Email, account.Password)
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	//temp := &Account{}

	//check for errors and duplicate emails
	//	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	//	if err != nil && err != gorm.ErrRecordNotFound {
	//		return u.Message(false, "Connection error. Please retry"), false
	//	}

	// пока заглушим
	//if ExistEmail(account.Email) {
	//	return u.Message(false, "Email address already in use by another user."), false
	//}

	return u.Message(false, "Requirement passed"), true
}

func (account *Account) Create() map[string]interface{} {

	if resp, ok := account.Validate(); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	//GetDB().Create(account)
	//account.ID = 16046
	CreateAccount(account)
	if account.ID == 0 {
		return u.Message(false, "Failed to create account, connection error or account exist")
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" //delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}

	err := GetAccount(account, email)
	if err == false {
		return u.Message(false, "Email address not found")
	}
	//err = GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	//if err != nil {
	//	if err == gorm.ErrRecordNotFound {
	//		return u.Message(false, "Email address not found")
	//	}
	//	return u.Message(false, "Connection error. Please retry")
	//}

	//err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	err2 := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err2 != nil && err2 == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	//Worked! Logged In
	//account.Password = ""

	//Create JWT token
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response
	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
	//if ValidateUser(email, password) {
	//		resp := u.Message(true, "Logged In")
	//return resp
	//} else {
	//resp := u.Message(false, "Not Logged")
	//return resp
	//}
}

func GetUser(u uint) *Account {

	acc := &Account{}
	//GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}

	acc.Password = ""
	return acc
}
