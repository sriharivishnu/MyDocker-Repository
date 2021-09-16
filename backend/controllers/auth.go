package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sriharivishnu/shopify-challenge/layers"
	"github.com/sriharivishnu/shopify-challenge/models"
	"github.com/sriharivishnu/shopify-challenge/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	UserService layers.UserLayer
}

func (controller *AuthController) SignUp(c *gin.Context) {
	var user models.User

	//input validation
	if errInputFormat := c.BindJSON(&user); errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}
	if errInputFormat := user.Validate(); errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	// hash user's password
	hashedPassword, errHash := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
	if errHash != nil {
		utils.RespondError(c, errHash, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// create the user
	user, errUserCreate := controller.UserService.Create(user.Username, user.Password)
	if errUserCreate != nil {
		utils.RespondError(c, errUserCreate, http.StatusInternalServerError)
		return
	}

	// generate token for user
	token, errToken := controller.UserService.CreateToken(user)
	if errToken != nil {
		utils.RespondError(c, errToken, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Signed up successfully", "token": token, "user": &user})

}

func (controller *AuthController) SignIn(c *gin.Context) {

	// Input validation
	var userPayload models.User
	errInputFormat := c.BindJSON(&userPayload)
	if errInputFormat != nil {
		utils.RespondError(c, errInputFormat, http.StatusBadRequest)
		return
	}

	// get user from db with matching username
	userDB, errGetUser := controller.UserService.GetByUsername(userPayload.Username)
	if errGetUser != nil {
		if errGetUser == sql.ErrNoRows {
			utils.RespondErrorString(c, "Unauthorized", http.StatusUnauthorized)
			return
		}
		utils.RespondError(c, errGetUser, http.StatusInternalServerError)
		return
	}

	// compare password
	if errHash := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(userPayload.Password)); errHash != nil {
		utils.RespondErrorString(c, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// create token for user
	token, errToken := controller.UserService.CreateToken(userDB)
	if errToken != nil {
		utils.RespondError(c, errToken, http.StatusInternalServerError)
		return
	}

	c.JSON(200, gin.H{"message": "Signed in successfully", "token": token})

}
