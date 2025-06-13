package mapper

import (
	pb "github.com/FlyKarlik/proto/spot_instrument_service/gen/spot_instrument_service/proto"
	"github.com/FlyKarlik/spotInstrumentService/internal/domain"
	"github.com/FlyKarlik/spotInstrumentService/pkg/proto_mapper"
)

func FromProtoViewMarketsRequest(pb *pb.ViewMarketsRequest) domain.ViewMarketsRequest {
	return domain.ViewMarketsRequest{
		UserRoles: FromProtoUserRoles(pb.UserRoles),
	}
}

func ToProtoViewMarketResponse(domain domain.ViewMarketsResponse) *pb.ViewMarketsResponse {
	return &pb.ViewMarketsResponse{
		Markets: ToProtoMarkets(domain.Markets),
	}
}

func ToProtoMarkets(markets []domain.Market) []*pb.Market {
	protoMarkets := make([]*pb.Market, 0, len(markets))
	for _, market := range markets {
		protoMarkets = append(protoMarkets, ToProtoMarket(market))
	}
	return protoMarkets
}

func ToProtoMarket(market domain.Market) *pb.Market {
	return &pb.Market{
		Id:           proto_mapper.ToIDProto(market.ID),
		Name:         proto_mapper.ToStringProto(market.Name),
		Enabled:      proto_mapper.ToBoolProto(market.Enabled),
		DeletedAt:    proto_mapper.ToTimestampProto(market.DeletedAt),
		AllowedRoles: ToProtoUserRoles(market.AllowedRoles),
	}
}

func FromProtoUserRole(protoRole pb.UserRole) domain.UserRoleEnum {
	switch protoRole {
	case pb.UserRole_USER_ROLE_TRADER:
		return domain.UserRoleEnumTrader
	case pb.UserRole_USER_ROLE_VIEWER:
		return domain.UserRoleEnumViewer
	case pb.UserRole_USER_ROLE_ADMIN:
		return domain.UserRoleEnumAdmin
	default:
		return domain.UserRoleEnumUnspecified
	}
}

func ToProtoUserRole(userRole domain.UserRoleEnum) pb.UserRole {
	switch userRole {
	case domain.UserRoleEnumTrader:
		return pb.UserRole_USER_ROLE_TRADER
	case domain.UserRoleEnumViewer:
		return pb.UserRole_USER_ROLE_VIEWER
	case domain.UserRoleEnumAdmin:
		return pb.UserRole_USER_ROLE_ADMIN
	default:
		return pb.UserRole_USER_ROLE_UNSPECIFIED
	}
}

func FromProtoUserRoles(protoRoles []pb.UserRole) domain.UserRolesEnum {
	userRoles := make(domain.UserRolesEnum, 0, len(protoRoles))
	for _, protoRole := range protoRoles {
		userRoles = append(userRoles, FromProtoUserRole(protoRole))
	}
	return userRoles
}

func ToProtoUserRoles(userRoles domain.UserRolesEnum) []pb.UserRole {
	protoRoles := make([]pb.UserRole, 0, len(userRoles))
	for _, userRole := range userRoles {
		protoRoles = append(protoRoles, ToProtoUserRole(userRole))
	}
	return protoRoles
}
