package response_service

import (
	"github.com/MamangRust/monolith-pointofsale-shared/domain/record"
	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"
)

type refreshTokenResponseMapper struct {
}

func NewRefreshTokenResponseMapper() *refreshTokenResponseMapper {
	return &refreshTokenResponseMapper{}
}

func (r *refreshTokenResponseMapper) ToRefreshTokenResponse(refresh *record.RefreshTokenRecord) *response.RefreshTokenResponse {
	return &response.RefreshTokenResponse{
		UserID:    refresh.UserID,
		Token:     refresh.Token,
		ExpiredAt: refresh.ExpiredAt,
		CreatedAt: refresh.CreatedAt,
		UpdatedAt: refresh.UpdatedAt,
	}
}

func (r *refreshTokenResponseMapper) ToRefreshTokenResponses(refreshs []*record.RefreshTokenRecord) []*response.RefreshTokenResponse {
	var responses []*response.RefreshTokenResponse

	for _, response := range refreshs {
		responses = append(responses, r.ToRefreshTokenResponse(response))
	}

	return responses
}
