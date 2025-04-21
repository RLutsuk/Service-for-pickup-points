package grpcPP

import (
	"context"

	pickupPointUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/usecase"
	ppProto "github.com/RLutsuk/Service-for-pickup-points/app/proto/pickuppoint"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	pickupPointUC pickupPointUC.UseCaseI
	ppProto.UnimplementedPPServiceServer
}

func NewGRPCServer(pickupPointUC pickupPointUC.UseCaseI) *GRPCServer {
	return &GRPCServer{pickupPointUC: pickupPointUC}
}

func (serv *GRPCServer) GetPickupPointList(context.Context, *ppProto.GetPickupPointRequest) (*ppProto.GetPickupPointResponse, error) {
	pickupPoints, err := serv.pickupPointUC.GetListOnlyPickupPoint()
	if err != nil {
		return nil, err
	}

	resPPs := make([]*ppProto.PickupPoint, 0)

	for _, PP := range pickupPoints {
		resPPs = append(resPPs, &ppProto.PickupPoint{Id: PP.ID, RegistrationDate: timestamppb.New(PP.RegistrationDate), City: PP.City})
	}

	return &ppProto.GetPickupPointResponse{PickupPoints: resPPs}, nil
}
