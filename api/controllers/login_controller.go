package controllers

import (
	"encoding/json"
	"github.com/nvthong2303/echo_crud_01/api/auth"
	"github.com/nvthong2303/echo_crud_01/api/models"
	"github.com/nvthong2303/echo_crud_01/api/responses"
	"github.com/nvthong2303/echo_crud_01/api/utils/formateror"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedErr := formateror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedErr)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {
	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", nil
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", nil
	}
	return auth.CreateToken(user.ID)
}
