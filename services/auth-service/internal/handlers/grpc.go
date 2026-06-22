package handlers

import (
	"context"
	"log"

	domain "famgo/auth-service/internal/domain"
	pb "famgo/auth-service/api/proto/v1"
)

// AuthServer implements auth.v1.AuthService
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service domain.AuthService
	logger  *log.Logger
}

// NewAuthServer creates gRPC server
func NewAuthServer(service domain.AuthService) *AuthServer {
	return &AuthServer{
		service: service,
		logger:  log.New(nil, "auth-server: ", log.LstdFlags),
	}
}

// Login authenticates user
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	s.logger.Printf("Login request for phone: %s", req.Phone)

	// Call domain service
	authReq := &domain.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
		DeviceID: req.DeviceId,
	}

	result, err := s.service.Login(ctx, authReq)
	if err != nil {
		s.logger.Printf("Login error: %v", err)
		return nil, err
	}

	// Map to proto response
	response := &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: &pb.UserProfile{
			Id:       result.User.ID,
			Phone:    result.User.Phone,
			Email:    result.User.Email,
			FullName: result.User.FullName,
		},
		Session: &pb.SessionData{
			Id:       result.Session.ID,
			DeviceId: result.Session.DeviceID,
		},
	}

	return response, nil
}

// Register creates new user
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	s.logger.Printf("Register request for phone: %s", req.Phone)

	domainReq := &domain.RegisterRequest{
		Phone:       req.Phone,
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		DeviceID:    req.DeviceId,
		Fingerprint: req.Fingerprint,
	}

	result, err := s.service.Register(ctx, domainReq)
	if err != nil {
		s.logger.Printf("Register error: %v", err)
		return nil, err
	}

	response := &pb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		User: &pb.UserProfile{
			Id:       result.User.ID,
			Phone:    result.User.Phone,
			Email:    result.User.Email,
			FullName: result.User.FullName,
		},
		Session: &pb.SessionData{
			Id:       result.Session.ID,
			DeviceId: result.Session.DeviceID,
		},
	}

	return response, nil
}

// VerifyToken validates JWT
func (s *AuthServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.TokenClaimsResponse, error) {
	claims, err := s.service.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	response := &pb.TokenClaimsResponse{
		UserId:      claims.UserID,
		Phone:       claims.Phone,
		Email:       claims.Email,
		Roles:       claims.Roles,
		MfaVerified: claims.MFAVerified,
		ExpiresAt:   claims.ExpiresAt.Unix(),
	}

	return response, nil
}

// RefreshToken generates new access token
func (s *AuthServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	result, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &pb.TokenResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}, nil
}

// CreateSession starts new session
func (s *AuthServer) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.SessionResponse, error) {
	domainReq := &domain.SessionRequest{
		DeviceID:    req.DeviceId,
		DeviceName:  req.DeviceName,
		IPAddress:   req.IpAddress,
		Fingerprint: req.Fingerprint,
		UserAgent:   req.UserAgent,
	}

	result, err := s.service.CreateSession(ctx, req.UserId, domainReq)
	if err != nil {
		return nil, err
	}

	return &pb.SessionResponse{
		Id:        result.ID,
		UserId:    result.UserID,
		DeviceId:  result.DeviceID,
		CreatedAt: result.CreatedAt.Unix(),
		ExpiresAt: result.ExpiresAt.Unix(),
	}, nil
}

// ValidateSession checks session validity
func (s *AuthServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	valid, err := s.service.ValidateSession(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateSessionResponse{
		Valid: valid,
	}, nil
}

// Logout terminates session
func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := s.service.Logout(ctx, req.SessionId)
	if err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{Success: true}, nil
}

// RevokeAllSessions invalidates all user sessions
func (s *AuthServer) RevokeAllSessions(ctx context.Context, req *pb.RevokeAllSessionsRequest) (*pb.RevokeAllSessionsResponse, error) {
	err := s.service.RevokeAllSessions(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.RevokeAllSessionsResponse{Success: true}, nil
}

// RegisterDevice stores device info
func (s *AuthServer) RegisterDevice(ctx context.Context, req *pb.RegisterDeviceRequest) (*pb.RegisterDeviceResponse, error) {
	domainReq := &domain.DeviceRegistrationRequest{
		DeviceID:    req.DeviceId,
		DeviceName:  req.DeviceName,
		DeviceType:  req.DeviceType,
		Fingerprint: req.Fingerprint,
		OSVersion:   req.OsVersion,
		AppVersion:  req.AppVersion,
	}

	err := s.service.RegisterDevice(ctx, req.UserId, domainReq)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterDeviceResponse{
		Success:  true,
		DeviceId: req.DeviceId,
	}, nil
}

// GenerateOTP creates OTP
func (s *AuthServer) GenerateOTP(ctx context.Context, req *pb.GenerateOTPRequest) (*pb.GenerateOTPResponse, error) {
	otp, err := s.service.GenerateOTP(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GenerateOTPResponse{
		Otp:       otp,
		ExpiresAt: 300, // 5 minutes
	}, nil
}

// VerifyOTP validates OTP
func (s *AuthServer) VerifyOTP(ctx context.Context, req *pb.VerifyOTPRequest) (*pb.VerifyOTPResponse, error) {
	valid, err := s.service.VerifyOTP(ctx, req.UserId, req.Otp)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyOTPResponse{Valid: valid}, nil
}

// EnableMFA activates MFA
func (s *AuthServer) EnableMFA(ctx context.Context, req *pb.EnableMFARequest) (*pb.EnableMFAResponse, error) {
	result, err := s.service.EnableMFA(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.EnableMFAResponse{
		Secret:      result.Secret,
		QrCode:      result.QRCode,
		BackupCodes: result.BackupCodes,
	}, nil
}

// VerifyMFA validates MFA token
func (s *AuthServer) VerifyMFA(ctx context.Context, req *pb.VerifyMFARequest) (*pb.VerifyMFAResponse, error) {
	valid, err := s.service.VerifyMFA(ctx, req.UserId, req.Token)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyMFAResponse{Valid: valid}, nil
}

// DisableMFA deactivates MFA
func (s *AuthServer) DisableMFA(ctx context.Context, req *pb.DisableMFARequest) (*pb.DisableMFAResponse, error) {
	// TODO: Implement MFA disable logic
	return &pb.DisableMFAResponse{Success: true}, nil
}

// GetRBAC retrieves user permissions
func (s *AuthServer) GetRBAC(ctx context.Context, req *pb.GetRBACRequest) (*pb.GetRBACResponse, error) {
	policy, err := s.service.GetRBAC(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetRBACResponse{
		Roles:       policy.Roles,
		Permissions: policy.Permissions,
	}, nil
}

// UpdateRBAC updates user roles
func (s *AuthServer) UpdateRBAC(ctx context.Context, req *pb.UpdateRBACRequest) (*pb.UpdateRBACResponse, error) {
	// TODO: Implement RBAC update logic
	return &pb.UpdateRBACResponse{Success: true}, nil
}

// Health checks service health
func (s *AuthServer) Health(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{
		Status:  "healthy",
		Version: "0.1.0",
	}, nil
}
