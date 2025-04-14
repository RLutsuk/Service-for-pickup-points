package usecase

import (
	"io"
	"log/slog"
	"testing"

	mockPickup "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository/mocks"
	mockProduct "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/repository/mocks"
	mockReception "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository/mocks"

	"github.com/RLutsuk/Service-for-pickup-points/app/models"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProd := mockProduct.NewMockRepositoryI(ctrl)
	mockRec := mockReception.NewMockRepositoryI(ctrl)
	mockPP := mockPickup.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	got := New(mockProd, mockRec, mockPP, logger)
	if got == nil {
		t.Errorf("New() returned nil")
	}
}

func Test_useCase_CreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProd := mockProduct.NewMockRepositoryI(ctrl)
	mockRec := mockReception.NewMockRepositoryI(ctrl)
	mockPP := mockPickup.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	uc := New(mockProd, mockRec, mockPP, logger)

	mockPP.EXPECT().
		GetPickupPointByID("pp-1").
		Return(nil)

	mockRec.EXPECT().
		GetOpenReceptionByPPID("pp-1").
		Return("rec-1", "2025-01-01T12:00:00Z", nil)

	mockProd.EXPECT().
		CreateProduct(gomock.Any()).
		Return(nil)

	mockPP.EXPECT().
		GetPickupPointByID("pp-missing").
		Return(models.ErrPickupPointDontExist)

	mockPP.EXPECT().
		GetPickupPointByID("pp-bad").
		Return(nil)

	mockRec.EXPECT().
		GetOpenReceptionByPPID("pp-bad").
		Return("", "", nil)

	mockPP.EXPECT().
		GetPickupPointByID("pp-type-error").
		Return(nil)

	mockRec.EXPECT().
		GetOpenReceptionByPPID("pp-type-error").
		Return("rec-type", "", nil)

	tests := []struct {
		name          string
		pickupPointID string
		typeProduct   string
		wantErr       bool
	}{
		{
			name:          "successful creation",
			pickupPointID: "pp-1",
			typeProduct:   "одежда",
			wantErr:       false,
		},
		{
			name:          "PP not found",
			pickupPointID: "pp-missing",
			typeProduct:   "обувь",
			wantErr:       true,
		},
		{
			name:          "there is no open reception",
			pickupPointID: "pp-bad",
			typeProduct:   "электроника",
			wantErr:       true,
		},
		{
			name:          "product type invalid",
			pickupPointID: "pp-type-error",
			typeProduct:   "едa",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uc.CreateProduct(tt.pickupPointID, tt.typeProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.TypeProduct != tt.typeProduct {
				t.Errorf("CreateProduct() product type = %v, want %v", got.TypeProduct, tt.typeProduct)
			}
		})
	}
}

func Test_useCase_DeleteLastProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProd := mockProduct.NewMockRepositoryI(ctrl)
	mockRec := mockReception.NewMockRepositoryI(ctrl)
	mockPP := mockPickup.NewMockRepositoryI(ctrl)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	uc := New(mockProd, mockRec, mockPP, logger)

	mockPP.EXPECT().
		GetPickupPointByID("pp-1").
		Return(nil)

	mockRec.EXPECT().
		GetOpenReceptionByPPID("pp-1").
		Return("rec-1", "2025-01-01T12:00:00Z", nil)

	mockProd.EXPECT().
		GetLastProdcutByDate("rec-1").
		Return("product-123", nil)

	mockProd.EXPECT().
		DeleteLastProduct("product-123").
		Return(nil)

	mockPP.EXPECT().
		GetPickupPointByID("pp-none").
		Return(nil)

	mockRec.EXPECT().
		GetOpenReceptionByPPID("pp-none").
		Return("", "", nil)

	mockPP.EXPECT().
		GetPickupPointByID("pp-error").
		Return(models.ErrBadData)

	tests := []struct {
		name          string
		pickupPointID string
		wantErr       bool
	}{
		{
			name:          "successful deletion",
			pickupPointID: "pp-1",
			wantErr:       false,
		},
		{
			name:          "there is no open reception",
			pickupPointID: "pp-none",
			wantErr:       true,
		},
		{
			name:          "PP not found",
			pickupPointID: "pp-error",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := uc.DeleteLastProduct(tt.pickupPointID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteLastProduct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
