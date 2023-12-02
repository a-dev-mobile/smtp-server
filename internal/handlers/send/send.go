package send

import (
	"context"
	"github.com/a-dev-mobile/smtp-server/internal/config"
	"github.com/a-dev-mobile/smtp-server/internal/utils"
	pb "github.com/a-dev-mobile/smtp-server/proto"
	"golang.org/x/exp/slog"
)


type SendServiceServer struct {
	pb.UnimplementedEmailSenderApiServer
	Config *config.Config
	Logger *slog.Logger
}

func NewServiceServer(cfg *config.Config, lg *slog.Logger) *SendServiceServer {
	return &SendServiceServer{
		Config: cfg,
		Logger: lg,
	}
}

func (s *SendServiceServer) SendEmail(ctx context.Context, req *pb.EmailSenderRequest) (*pb.EmailSenderResponse, error) {
	s.Logger.Info("SendEmail called", "recipient", req.GetRecipientEmail())

	if !utils.ValidateEmail(req.GetRecipientEmail()) {
		return s.handleError("Invalid email format", req.GetRecipientEmail())
	}

	emailConfigs := utils.NewEmailConfigs(s.Config.SMTPProviders, s.Logger)
	err := utils.SendEmail(emailConfigs, req.GetFromName(), req.GetFromEmail(), req.GetRecipientEmail(), req.GetSubject(), req.GetBody())
	if err != nil {
		return s.handleError(err.Error(), req.GetRecipientEmail())
	}

	s.Logger.Info("Email sent successfully", "recipient", req.GetRecipientEmail())
	return &pb.EmailSenderResponse{Success: true, Message: "The email was successfully delivered"}, nil
}

func (s *SendServiceServer) handleError(msg, email string) (*pb.EmailSenderResponse, error) {
	s.Logger.Error(msg, "email", email)
	return &pb.EmailSenderResponse{Success: false, Message: msg}, nil
}
