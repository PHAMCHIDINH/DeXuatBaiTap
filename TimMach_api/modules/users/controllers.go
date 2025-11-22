package users

import (
	db "chidinh/db/sqlc"
	"chidinh/modules/auth"
	"chidinh/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Handler gom các dependency cho module users.
type Handler struct {
	Queries *db.Queries
	Tokens  auth.TokenService
}

func NewHandler(queries *db.Queries, tokens auth.TokenService) *Handler {
	return &Handler{
		Queries: queries,
		Tokens:  tokens,
	}
}

// POST /users/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Email == "" || len(req.Password) < 6 {
		utils.RespondError(c, http.StatusBadRequest, "email/password is invalid")
		return
	}

	_, err := h.Queries.GetUserByEmail(c, req.Email)
	if err == nil {
		utils.RespondError(c, http.StatusConflict, "email already exists")
		return
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("register: cannot check email %s: %v", req.Email, err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot check email")
		return
	}

	seq, err := h.Queries.NextUserSeq(c)
	if err != nil {
		log.Printf("register: cannot get user seq: %v", err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot create user")
		return
	}
	newID := utils.FormatUserID(seq, time.Now())
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("register: cannot hash password for %s: %v", req.Email, err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot hash password")
		return
	}

	user, err := h.Queries.CreateUser(c, db.CreateUserParams{
		ID:           newID,
		Email:        req.Email,
		PasswordHash: string(hashed),
	})
	if err != nil {
		log.Printf("register: cannot create user %s: %v", req.Email, err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot create user")
		return
	}
	token, err := h.Tokens.GenerateToken(c, user)
	if err != nil {
		log.Printf("register: cannot create token for %s: %v", req.Email, err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot create token")
		return
	}
	c.JSON(http.StatusCreated, RegisterResponse{
		User:  toUserResponse(toUserDomain(user)),
		Token: token,
	})
}

// POST /users/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.Queries.GetUserByEmail(c, req.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		log.Printf("login: cannot fetch user %s: %v", req.Email, err)
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch user")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := h.Tokens.GenerateToken(c, user)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot create token")
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: token,
		User:        toUserResponse(toUserDomain(user)),
	})
}

// GET /users/me
func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	user, err := h.Queries.GetUserByID(c, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch user")
		return
	}

	c.JSON(http.StatusOK, toUserResponse(toUserDomain(user)))
}
