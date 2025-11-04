package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	usermodel "quepc/api/internal/user/model"
	"quepc/api/internal/user/store"
)

type Handler struct {
	DB *gorm.DB
	US *store.Store
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db, US: store.New(db)}
}

var allowedRoles = map[string]bool{
	"ADMIN":  true,
	"EDITOR": true,
}

type createUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"` // "ADMIN" | "EDITOR"
}

type createUserRes struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

func (h *Handler) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || len(req.Password) < 8 {
		http.Error(w, "email inválido o password muy corto (mín. 8)", http.StatusBadRequest)
		return
	}
	if req.Role == "" {
		req.Role = "VIEWER"
	}
	if !allowedRoles[strings.ToUpper(req.Role)] {
		http.Error(w, "role inválido (ADMIN|EDITOR)", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "no se pudo hashear password", http.StatusInternalServerError)
		return
	}

	u := &usermodel.User{
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
		Role:     strings.ToUpper(req.Role),
	}
	if err := h.US.Create(u); err != nil {
		http.Error(w, "email ya usado o error al crear", http.StatusConflict)
		return
	}

	writeJSON(w, createUserRes{
		ID: u.ID, Email: u.Email, Name: u.Name, Role: u.Role,
	}, http.StatusCreated)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type loginRes struct {
	Token   string `json:"token"`
	Expires string `json:"expires_at"`
	User    struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	} `json:"user"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	u, err := h.US.FindByEmail(req.Email)
	if err != nil {
		http.Error(w, "credenciales inválidas", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		http.Error(w, "credenciales inválidas", http.StatusUnauthorized)
		return
	}
	tok, exp, err := SignJWT(u.ID, u.Email, u.Role)
	if err != nil {
		http.Error(w, "no se pudo firmar token", http.StatusInternalServerError)
		return
	}
	var res loginRes
	res.Token = tok
	res.Expires = exp.UTC().Format(timeRFC3339)
	res.User.ID, res.User.Email, res.User.Name, res.User.Role = u.ID, u.Email, u.Name, u.Role
	writeJSON(w, res, http.StatusOK)
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	uid, ok := UserIDFromContext(r)
	if !ok {
		http.Error(w, "no auth context", http.StatusUnauthorized)
		return
	}
	writeJSON(w, map[string]any{"user_id": uid}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if !strings.HasPrefix(ah, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		claims, err := ParseJWT(strings.TrimPrefix(ah, "Bearer "))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		if claims.Role != "ADMIN" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Crea un usuario ADMIN si no existe ninguno en la base de datos
func CrearSuperAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&usermodel.User{}).Where("role = ?", "ADMIN").Count(&count).Error; err != nil {
		return fmt.Errorf("error verificando admin: %w", err)
	}
	if count > 0 {
		fmt.Println("Ya existe al menos un ADMIN, no se crea ninguno nuevo.")
		return nil
	}

	email := strings.ToLower(os.Getenv("ADMIN_EMAIL"))
	password := os.Getenv("ADMIN_PASSWORD")
	name := os.Getenv("ADMIN_NAME")
	if name == "" {
		name = "Super Admin"
	}

	if email == "" || password == "" {
		return fmt.Errorf("faltan variables de entorno ADMIN_EMAIL y ADMIN_PASSWORD")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generando hash: %w", err)
	}

	admin := &usermodel.User{
		Email:    email,
		Password: string(hash),
		Name:     name,
		Role:     "ADMIN",
	}

	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("error creando superadmin: %w", err)
	}

	fmt.Println("Usuario ADMIN creado:")
	fmt.Println("   Email:", admin.Email)
	fmt.Println("   Password:", password)
	return nil
}
