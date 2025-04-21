package usecase

import (
	"errors"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	mocks "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository/mocks"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	got := New(mockRepo, logger)
	if got == nil {
		t.Errorf("New() returned nil")
	}
}

func Test_useCase_CreatePickupPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)

	validPickup := &models.PickupPoint{ID: "1", City: "Москва"}
	invalidPickup := &models.PickupPoint{ID: "2", City: "Лондон"}

	mockRepo.EXPECT().
		CreatePickupPoint(validPickup).
		Return(nil).
		Times(1)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepo, logger)

	tests := []struct {
		name    string
		point   *models.PickupPoint
		wantErr bool
	}{
		{
			name:    "valid city",
			point:   validPickup,
			wantErr: false,
		},
		{
			name:    "invalid city",
			point:   invalidPickup,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.CreatePickupPoint(tt.point)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePickupPoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_useCase_GetAllPickupPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)

	now := time.Now().Format(time.RFC3339)

	mockRepo.EXPECT().
		GetAllPickupPoint("", "", 0, 5).
		Return([]*models.PickupPoint{{ID: "1"}}, nil).
		Times(1)

	mockRepo.EXPECT().
		GetAllPickupPoint(now, now, 0, 5).
		Return([]*models.PickupPoint{{ID: "2"}}, nil).
		Times(1)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepo, logger)

	tests := []struct {
		name      string
		startDate string
		endDate   string
		page      string
		limit     string
		want      []*models.PickupPoint
		wantErr   bool
	}{
		{
			name:    "default params",
			page:    "1",
			limit:   "5",
			want:    []*models.PickupPoint{{ID: "1"}},
			wantErr: false,
		},
		{
			name:      "with dates and incorrect nums",
			startDate: now,
			endDate:   now,
			page:      "-1",
			limit:     "-5",
			want:      []*models.PickupPoint{{ID: "2"}},
			wantErr:   false,
		},
		{
			name:    "invalid limit",
			page:    "1",
			limit:   "not-a-number",
			wantErr: true,
		},
		{
			name:    "invalid page",
			page:    "abc",
			limit:   "5",
			wantErr: true,
		},
		{
			name:      "invalid date format",
			startDate: "2025-13-01",
			page:      "1",
			limit:     "5",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.GetAllPickupPoint(tt.startDate, tt.endDate, tt.page, tt.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPickupPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllPickupPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_useCase_GetListOnlyPickupPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryI(ctrl)

	mockRepo.EXPECT().GetListOnlyPickupPoint().Times(1).
		Return([]*models.PickupPoint{{ID: "1"}}, nil)

	mockRepo.EXPECT().GetListOnlyPickupPoint().Times(1).
		Return(nil, errors.New("database error"))

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepo, logger)

	tests := []struct {
		name    string
		want    []*models.PickupPoint
		wantErr bool
	}{
		{
			name:    "successful request",
			want:    []*models.PickupPoint{{ID: "1"}},
			wantErr: false,
		},
		{
			name:    "database error",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.GetListOnlyPickupPoint()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePickupPoint() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllPickupPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
