package user

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/bhlox/ecom/internal/configs"
	"github.com/bhlox/ecom/internal/db"
	"github.com/bhlox/ecom/internal/response"
	"github.com/bhlox/ecom/internal/services/auth"
	"github.com/bhlox/ecom/internal/types"
	"github.com/bhlox/ecom/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserQueries interface {
	GetUser(ctx context.Context, arg db.GetUserParams) (db.User, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetAllUsers(ctx context.Context) ([]db.User, error)
	DeleteUser(ctx context.Context, arg db.DeleteUserParams) error
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var data types.LoginPayload
	if err := utils.ParseJSONReq(r, &data); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	user, errs := h.queries.GetUser(r.Context(), db.GetUserParams{Email: data.Email})
	if errs != nil {
		response.Error(w, http.StatusNotFound, errs.Error())
		return
	}
	if err := verifyPW(data.Password, user.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := auth.GenerateJWT(configs.Envs.JWTSECRET, user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var data types.RegisterPayload
	if err := utils.ParseJSONReq(r, &data); err != nil {
		fmt.Println("err processing JSON")
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.Validate.Struct(data); err != nil {
		// fmt.Println(err.Error())
		ValidationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			fmt.Println("err is not that type")
			response.Error(w, 500, "wtf")
			return
		}
		var errors = make(map[string]string)
		for _, validationError := range ValidationErrors {
			errors[validationError.Field()] = validationError.Error()
		}
		fmt.Println("failed on validating payload")
		response.JSON(w, http.StatusNotAcceptable, errors)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	// error would occur if no user is found using the searched params
	foundUser, err := h.queries.GetUser(r.Context(), db.GetUserParams{Email: data.Email})
	if !foundUser.CreatedAt.IsZero() && err == nil {
		// fmt.Printf("An email of %v already exists\n", data.Email)
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("An email of %v already exists\n", data.Email))
		return
	}
	// fmt.Println(foundUser)

	hashedPassword, err := utils.HashString(data.Password)
	if err != nil {
		response.Error(w, 500, "error hashing password")
		return
	}
	user, errs := h.queries.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		FirstName: data.Firstname,
		LastName:  data.Lastname,
		Email:     data.Email,
		Password:  hashedPassword,
	})
	if errs != nil {
		fmt.Println("failed on creating user")
		response.Error(w, 401, errs.Error())
		return
	}
	response.JSON(w, http.StatusCreated, user)
}

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	users, err := h.queries.GetAllUsers(r.Context())
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, users)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, 400, "no params found")
	}
	userUUID, err := uuid.Parse(id)
	if err != nil {
		response.Error(w, 400, "invalid ID format")
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	// Use the UUID in the DeleteUserParams
	if err := h.queries.DeleteUser(r.Context(), db.DeleteUserParams{ID: userUUID}); err != nil {
		response.Error(w, 500, "error deleting user")
		return
	}
	response.JSON(w, 200, "deleted")
}

func verifyPW(payloadPW, storedPW string) error {
	decodedPW, err := base64.StdEncoding.DecodeString(storedPW)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(decodedPW, []byte(payloadPW)); err != nil {
		return err
	}
	return nil
}
