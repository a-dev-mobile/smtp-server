package smtpserver

// import (
// 	"bufio"
// 	"crypto/tls"
// 	"fmt"
// 	"io"
// 	"smtp-server/internal/config" // Импортируем ваш пакет конфигурации

// 	"net"
// 	"net/smtp"
// 	"smtp-server/internal/environment"
// 	"strings"
// 	"time"

// 	"golang.org/x/exp/slog"
// )

// // SMTPServer представляет ваш SMTP сервер
// type SMTPServer struct {
// 	config      *config.SMTPServerConfig
// 	Logger      *slog.Logger
// 	Environment environment.Environment
// }

// // NewSMTPServer создает новый экземпляр SMTPServer с данными конфигурации и логгером
// func NewSMTPServer(cfg *config.SMTPServerConfig, lg *slog.Logger, env environment.Environment) *SMTPServer {
// 	return &SMTPServer{
// 		config:      cfg,
// 		Logger:      lg,
// 		Environment: env,
// 	}
// }

// // Start запускает SMTP сервер
// func (s *SMTPServer) Start() {
// 	var listener net.Listener
// 	var err error
// 	s.Logger.Info("SMTP server is starting...", "port", s.config.Port)

// 	if s.Environment == environment.Dev {
// 		s.Logger.Info("Running in development mode, no TLS")
// 		listener, err = net.Listen("tcp", ":"+s.config.Port)

// 		if err != nil {
// 			s.Logger.Error("Failed to start TCP listener:", "error", err)
// 			return
// 		}
// 	} else {
// 		cert, err := tls.LoadX509KeyPair(s.config.TLSCertPath, s.config.TLSKeyPath)
// 		if err != nil {
// 			s.Logger.Error("Failed to load", "key pair", err)
// 			return
// 		}
// 		tlsConfig := &tls.Config{
// 			Certificates: []tls.Certificate{cert},
// 			MinVersion:   tls.VersionTLS12}

// 		listener, err = tls.Listen("tcp", ":"+s.config.Port, tlsConfig)
// 		if err != nil {
// 			s.Logger.Error("Failed to start TLS listener:", "error", err)
// 			return
// 		}
// 	}

// 	defer listener.Close()
// 	s.Logger.Info("SMTP server is running", "port", s.config.Port)


//     if s.Environment == environment.Dev {
//         s.Logger.Info("Sending test email on startup")
//         err := SendTestEmail(s)
//         if err != nil {
//             s.Logger.Error("Failed to send test email", "error", err)
//         } else {
//             s.Logger.Info("Test email sent successfully")
//         }
//     }



// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			s.Logger.Error("Failed to accept connection:", "error", err)
// 			continue
// 		}
// 		s.Logger.Info("Connection accepted", "remote_addr", conn.RemoteAddr().String())

// 		go handleConnection(conn, s)
// 	}

// }

// // SendTestEmail sends a test email
// func SendTestEmail(s *SMTPServer) error {
// 	from := "noreply@wayofdt.com"
// 	to := "a.dev.mobile@gmail.com"
// 	subject := "Test Email"
// 	body := "This is a test email sent from SMTP Server."

// 	return SendEmail(from, to, subject, body, s)
// }

// // handleConnection обрабатывает входящие SMTP соединения
// func handleConnection(conn net.Conn, s *SMTPServer) {

// 	defer conn.Close()

// 	s.Logger.Info("Connection established", "remote_addr", conn.RemoteAddr().String())

// 	// Установка таймаутов для соединения
// 	if err := conn.SetReadDeadline(time.Now().Add(1 * time.Minute)); err != nil {
// 		s.Logger.Error("Failed to set read deadline", "error", err)
// 	}
// 	if err := conn.SetWriteDeadline(time.Now().Add(1 * time.Minute)); err != nil {
// 		s.Logger.Error("Failed to set write deadline", "error", err)
// 	}
// 	reader := bufio.NewReader(conn)
// 	writer := bufio.NewWriter(conn)

// 	// Отправляем приветственное сообщение
// 	writeResponse(writer, "220 Welcome to My SMTP Server", s)

// 	var from, to, subject, body string
// 	readingData := false

// 	for {
// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			handleReadError(err, s)
// 			break
// 		}
// 		// Ignore empty commands
// 		if strings.TrimSpace(line) == "" {
// 			continue
// 		}
// 		// Логирование принятой команды
// 		s.Logger.Info("Command received", "command", line)

// 		if readingData {
// 			if line == "." {
// 				readingData = false
// 				subject, body = parseData(subject, body)
// 				handleSendEmail(writer, from, to, subject, body, s)
// 				continue
// 			}
// 			body += line
// 		} else {
// 			command := strings.ToUpper(strings.TrimSpace(line))
// 			switch {
// 			case strings.HasPrefix(command, "EHLO"):
// 				writeResponse(writer, "250-Hello "+strings.TrimSpace(strings.TrimPrefix(command, "EHLO")), s)
// 				if s.Environment != environment.Dev {
// 					writeResponse(writer, "250-STARTTLS", s)
// 				}
// 				writeResponse(writer, "250 OK", s)
// 			case strings.HasPrefix(command, "MAIL FROM:"):
// 				from = extractEmail(command)
// 				writeResponse(writer, "250 OK", s)
// 			case strings.HasPrefix(command, "RCPT TO:"):
// 				to = extractEmail(command)
// 				writeResponse(writer, "250 OK", s)
// 			case strings.HasPrefix(command, "DATA"):
// 				writeResponse(writer, "354 Start mail input", s)
// 				readingData = true
// 			case strings.HasPrefix(command, "QUIT"):
// 				writeResponse(writer, "221 Bye", s)
// 				return
// 			default:
// 				writeResponse(writer, "502 Command not implemented", s)
// 			}
// 		}
// 	}
// }

// // handleReadError обрабатывает ошибки чтения из соединения
// func handleReadError(err error, s *SMTPServer) {
// 	if err != io.EOF {
// 		s.Logger.Error("Error reading from connection", "error", err)
// 	} else {
// 		s.Logger.Info("Connection closed by client")
// 	}
// }

// // handleSendEmail обрабатывает отправку email
// func handleSendEmail(writer *bufio.Writer, from, to, subject, body string, s *SMTPServer) {
// 	s.Logger.Info("Preparing to send email", "from", from, "to", to)
// 	err := SendEmail(from, to, subject, body, s)
// 	if err != nil {
// 		s.Logger.Error("Error in sending email", "error", err)
// 		writeResponse(writer, fmt.Sprintf("550 Error sending email: %v", err), s)
// 		return
// 	}
// 	writeResponse(writer, "250 Message sent successfully", s)
// }

// // parseData разбирает данные и извлекает тему сообщения
// func parseData(subject, body string) (string, string) {
// 	// Тут логика разбора темы сообщения из тела, если это необходимо
// 	// Пример:
// 	if strings.Contains(body, "Subject:") {
// 		subject = strings.Split(body, "Subject:")[1]
// 		subject = strings.Split(subject, "\n")[0]
// 	}
// 	return subject, body
// }

// // writeResponse отправляет ответ клиенту
// func writeResponse(writer *bufio.Writer, response string, s *SMTPServer) {
// 	s.Logger.Info("Sending response", "response", response)
// 	writer.WriteString(response + "\r\n")
// 	if err := writer.Flush(); err != nil {
// 		s.Logger.Error("Error while flushing writer", "error", err)
// 	}
// }

// // extractEmail извлекает email адрес из команды
// func extractEmail(command string) string {
// 	// Предполагается, что команда выглядит как "MAIL FROM:<email>" или "RCPT TO:<email>"
// 	start := strings.Index(command, "<") + 1
// 	end := strings.Index(command, ">")
// 	if start < 1 || end < 0 || end <= start {
// 		return ""
// 	}
// 	return command[start:end]
// }
// func SendEmail(from, to, subject, body string, s *SMTPServer) error {
// 	// Формируем заголовки и тело сообщения
// 	message := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body))

// 	// Использование доменного имени и порта из конфигурации
// 	smtpServer := s.config.Domain
// 	smtpPort := s.config.Port

// 	// Здесь мы отправляем сообщение непосредственно, без аутентификации, так как это наш сервер
// 	client, err := smtp.Dial(fmt.Sprintf("%s:%s", smtpServer, smtpPort))
// 	if err != nil {
// 		return fmt.Errorf("failed to connect to SMTP server: %v", err)
// 	}
// 	defer client.Close()

// 	if err = client.Mail(from); err != nil {
// 		return fmt.Errorf("failed to set sender: %v", err)
// 	}

// 	if err = client.Rcpt(to); err != nil {
// 		return fmt.Errorf("failed to set recipient: %v", err)
// 	}

// 	wc, err := client.Data()
// 	if err != nil {
// 		return fmt.Errorf("failed to create data writer: %v", err)
// 	}
// 	defer wc.Close()

// 	if _, err = wc.Write(message); err != nil {
// 		return fmt.Errorf("failed to write message: %v", err)
// 	}

// 	return nil
// }
