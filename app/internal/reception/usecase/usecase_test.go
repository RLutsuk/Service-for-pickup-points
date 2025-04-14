package usecase

import (
	"errors"
	"io"
	"log/slog"
	"reflect"
	"testing"
	"time"

	mocksPP "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository/mocks"
	mocksRec "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository/mocks"
	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoPP := mocksPP.NewMockRepositoryI(ctrl)
	mockRepoRec := mocksRec.NewMockRepositoryI(ctrl)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	got := New(mockRepoRec, mockRepoPP, logger)
	if got == nil {
		t.Errorf("New() returned nil")
	}
}

func Test_useCase_CreateReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoPP := mocksPP.NewMockRepositoryI(ctrl)
	mockRepoRec := mocksRec.NewMockRepositoryI(ctrl)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepoRec, mockRepoPP, logger)

	receptionSuccess := &models.Reception{ID: "1", PickupPointID: "pp-success"}
	receptionNoPVZ := &models.Reception{ID: "2", PickupPointID: "pp-not-found"}
	receptionOpenExists := &models.Reception{ID: "3", PickupPointID: "pp-already-open"}
	receptionErrOnGetOpen := &models.Reception{ID: "4", PickupPointID: "pp-get-open-error"}
	receptionErrOnCreate := &models.Reception{ID: "5", PickupPointID: "pp-create-error"}
	receptionBadPVZ := &models.Reception{ID: "6", PickupPointID: "pp-internal-error"}

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-success").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-success").
		Return("", "", nil)

	mockRepoRec.EXPECT().
		CreateReception(receptionSuccess).
		Return(nil)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-not-found").
		Return(models.ErrPickupPointDontExist)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-already-open").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-already-open").
		Return("open-id", "2025-01-01T00:00:00Z", nil)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-get-open-error").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-get-open-error").
		Return("", "", errors.New("unexpected db error"))

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-create-error").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-create-error").
		Return("", "", nil)

	mockRepoRec.EXPECT().
		CreateReception(receptionErrOnCreate).
		Return(errors.New("insert failed"))

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-internal-error").
		Return(errors.New("some db error"))

	tests := []struct {
		name    string
		input   *models.Reception
		wantErr bool
	}{
		{
			name:    "successful creation",
			input:   receptionSuccess,
			wantErr: false,
		},
		{
			name:    "PP not found",
			input:   receptionNoPVZ,
			wantErr: true,
		},
		{
			name:    "there is already an open acceptance",
			input:   receptionOpenExists,
			wantErr: true,
		},
		{
			name:    "unsuccessful GetOpenReceptionByPPID",
			input:   receptionErrOnGetOpen,
			wantErr: true,
		},
		{
			name:    "unsuccessful CreateReception",
			input:   receptionErrOnCreate,
			wantErr: true,
		},
		{
			name:    "internal error",
			input:   receptionBadPVZ,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.CreateReception(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("useCase.CreateComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_useCase_CloseReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoPP := mocksPP.NewMockRepositoryI(ctrl)
	mockRepoRec := mocksRec.NewMockRepositoryI(ctrl)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	uc := New(mockRepoRec, mockRepoPP, logger)

	const validPVZ = "pp-1"
	const openID = "rec-123"
	const validTimeStr = "2025-04-13T12:00:00Z"
	parsedTime, _ := time.Parse(time.RFC3339, validTimeStr)

	mockRepoPP.EXPECT().
		GetPickupPointByID(validPVZ).
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID(validPVZ).
		Return(openID, validTimeStr, nil)

	mockRepoRec.EXPECT().
		CloseReception(openID).
		Return(nil)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-not-found").
		Return(models.ErrPickupPointDontExist)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-bad-data").
		Return(errors.New("some db error"))

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-no-open").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-no-open").
		Return("", "", nil)

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-get-error").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-get-error").
		Return("", "", errors.New("get error"))

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-close-error").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-close-error").
		Return(openID, validTimeStr, nil)

	mockRepoRec.EXPECT().
		CloseReception(openID).
		Return(errors.New("close error"))

	mockRepoPP.EXPECT().
		GetPickupPointByID("pp-parse-error").
		Return(nil)

	mockRepoRec.EXPECT().
		GetOpenReceptionByPPID("pp-parse-error").
		Return(openID, "bad-time-format", nil)

	mockRepoRec.EXPECT().
		CloseReception(openID).
		Return(nil)

	tests := []struct {
		name    string
		inputID string
		want    *models.Reception
		wantErr bool
	}{
		{
			name:    "successful closure",
			inputID: validPVZ,
			want: &models.Reception{
				ID:            openID,
				DateTime:      parsedTime,
				PickupPointID: validPVZ,
				Status:        "close",
			},
			wantErr: false,
		},
		{
			name:    "ПВЗ не найден",
			inputID: "pp-not-found",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ошибка при получении ПВЗ",
			inputID: "pp-bad-data",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "нет открытой приёмки",
			inputID: "pp-no-open",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ошибка при GetOpenReceptionByPPID",
			inputID: "pp-get-error",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ошибка при CloseReception",
			inputID: "pp-close-error",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ошибка при time.Parse",
			inputID: "pp-parse-error",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.CloseReception(tt.inputID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloseReception() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloseReception() = %v, want %v", got, tt.want)
			}
		})
	}
}
