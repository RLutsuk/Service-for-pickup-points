package usecase

import (
	"io"
	"log/slog"
	"testing"
	"time"

	mocks "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/repository/mocks"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	got := New(mockUserRepo, logger)
	if got == nil {
		t.Errorf("New() returned nil")
	}
}

func Test_useCase_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepo, logger)

	validUser := &models.User{
		Email:    "test@example.com",
		Password: "securepass",
		Role:     "employee",
	}

	mockRepo.EXPECT().
		ChekUserByEmail("test@example.com").
		Return(nil)

	mockRepo.EXPECT().
		CreateUser(gomock.Any()).
		DoAndReturn(func(u *models.User) (*models.User, error) {
			u.ID = "generated-id"
			return u, nil
		})

	tests := []struct {
		name    string
		user    *models.User
		mock    func()
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    validUser,
			mock:    func() {},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: &models.User{Email: "invalid", Password: "pass", Role: "employee"},
			mock: func() {
				mockRepo.EXPECT().
					ChekUserByEmail("invalid").
					Return(nil)
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			user: &models.User{Email: "valid@mail.com", Password: "pass", Role: "wrong"},
			mock: func() {
				mockRepo.EXPECT().
					ChekUserByEmail("valid@mail.com").
					Return(nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := uc.CreateUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_useCase_AuthUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepo, logger)

	rawPassword := "testpass"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	user := &models.User{
		Email:    "auth@mail.com",
		Password: rawPassword,
		Role:     "employee",
	}

	mockRepo.EXPECT().
		GetUserByEmail(user).
		DoAndReturn(func(u *models.User) error {
			u.ID = "u-123"
			u.Password = string(hashed)
			u.Role = "employee"
			return nil
		})

	token, err := uc.AuthUser(user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Error("expected token, got empty string")
	}
}

func Test_useCase_TestUser(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(nil, logger)

	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{"valid employee", "employee", false},
		{"valid moderator", "moderator", false},
		{"invalid role", "admin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := uc.TestUser(tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("TestUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Error("expected token, got empty")
			}
		})
	}
}

func TestParsetoken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId:   "test-id",
		UserRole: "moderator",
	})
	signed, _ := token.SignedString([]byte(signingKey))

	tests := []struct {
		name        string
		accessToken string
		wantID      string
		wantRole    string
		wantErr     bool
	}{
		{
			name:        "valid token",
			accessToken: signed,
			wantID:      "test-id",
			wantRole:    "moderator",
			wantErr:     false,
		},
		{
			name:        "invalid token",
			accessToken: "not-a-token",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, role, err := Parsetoken(tt.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parsetoken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.wantID {
				t.Errorf("Parsetoken() id = %v, want %v", id, tt.wantID)
			}
			if role != tt.wantRole {
				t.Errorf("Parsetoken() role = %v, want %v", role, tt.wantRole)
			}
		})
	}
}
