package services

import (
	"errors"
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
	logRepo  *repositories.AuditLogRepository
	cfg      config.Config
}

func NewAuthService(ur *repositories.UserRepository, lr *repositories.AuditLogRepository, cfg config.Config) *AuthService {
	return &AuthService{userRepo: ur, logRepo: lr, cfg: cfg}
}

// Login 校验口令并签发 JWT。
func (s *AuthService) Login(c *gin.Context, req types.LoginRequest) (*types.LoginView, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.NewAppError(401, constants.AuthInvalid, constants.AuthInvalidMessage)
	}
	if err != nil {
		return nil, err
	}
	if user.PasswordHash != utils.HashPassword(req.Password, s.cfg.JWTSecret) {
		return nil, utils.NewAppError(401, constants.AuthInvalid, constants.AuthInvalidMessage)
	}
	ttl, _ := time.ParseDuration(s.cfg.TokenTTLHour + "h")
	identity := utils.AuthIdentity{Username: user.Username, Role: user.Role, TeamCode: user.TeamCode, Name: user.DisplayName}
	token := utils.IssueJWT(identity, s.cfg.JWTSecret, ttl)
	return &types.LoginView{
		Token: token, Username: user.Username, Role: user.Role,
		RoleText: constants.StatusText["Role"][user.Role], TeamCode: user.TeamCode, DisplayName: user.DisplayName,
	}, nil
}

func (s *AuthService) ListLogs(limit int) ([]types.AuditLogView, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.logRepo.List(limit)
	if err != nil {
		return nil, err
	}
	views := make([]types.AuditLogView, 0, len(rows))
	for i := range rows {
		views = append(views, constructors.BuildAuditLogView(rows[i]))
	}
	return views, nil
}
