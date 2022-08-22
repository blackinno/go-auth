package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"auth.jwt.api/database"
	"auth.jwt.api/models"
	"auth.jwt.api/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleConf *oauth2.Config

type registerRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
}

func init() {
	googleConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  "http://localhost:8080/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func Register(c *gin.Context) {
	var payload registerRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist models.User

	database.DB.Where("username = ? ", payload.Username).First(&userExist)

	if userExist.ID > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User is already existing."})
		return
	}

	passwordEncrypt, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{Username: payload.Username, Password: string(passwordEncrypt), Fullname: payload.Fullname, Provider: "local"}

	database.DB.Create(&user)

	c.JSON(http.StatusOK, user)
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	status      string ``
	accessToken string ``
}

const (
	SUCCESS = "SUCCESS"
	FAILED  = "FAILED"
)

func Login(c *gin.Context) {
	var payload loginRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	record := database.DB.Where("username = ? ", payload.Username).First(&user)

	if record.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Username/Password is wrong or not exist"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Username/Password is wrong or not exist"})
		return
	}

	token, err := util.GenerateJWT(user.Username, int(user.ID))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": SUCCESS, "accessToken": token})
}

func GoogleAuth(c *gin.Context) {
	redirectUrl := googleConf.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

type googleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Locale        string `json:"locale"`
}

func GoogleAuthCallback(c *gin.Context) {
	state := c.Query("state")

	if state != "state" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	code := c.Query("code")

	token, err := googleConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	client := googleConf.Client(oauth2.NoContext, token)

	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user googleUser
	json.Unmarshal(content, &user)

	c.JSON(http.StatusOK, user)
}
