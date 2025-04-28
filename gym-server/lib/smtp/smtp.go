package smtp

import (
	"fmt"
	"health/internal/config"

	"gopkg.in/mail.v2"
)

var textMessage = "Уважаемый пользователь!\n\n" +
	"Сервис производственная гимнастика запрашивает подтверждение о том, что email %v — действительно ваш.\n\n" +
	"Токен подтверждения: %v\n\n" +
	"С уважением, сервис Производственная гимнастика.\n\n" +
	"Если вы не запрашивали подтверждение, просто проигнорируйте это сообщение."

func SendMail(smtpConfig config.SMTPConfig, authToken, userEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", smtpConfig.EmailSender)
	m.SetHeader("To", userEmail)
	m.SetAddressHeader("Cc", userEmail, "")
	m.SetHeader("Subject", "Gymanstic: Подтвердждения email")

	textMessage := fmt.Sprintf(textMessage, userEmail, authToken)
	m.SetBody("text/plain", textMessage)

	// Тут должны быть разные методы (verify и reset)
	verifyURL := fmt.Sprintf("https://your-app.com/verify?token=%s", authToken)
	htmlBody := fmt.Sprintf(`
	<html>
		<body>
			<p>Уважаемый пользователь!</p>
			<p>Сервис производственная гимнастика запрашивает подтверждение, что email <b>%s</b> действительно ваш.</p>
			<p>
				<a href="%s" style="
					display: inline-block;
					padding: 12px 24px;
					font-size: 16px;
					color: #ffffff;
					background-color: #007BFF;
					text-decoration: none;
					border-radius: 6px;
				">Подтвердить Email</a>
			</p>
			<p>Или скопируйте ссылку в браузер:</p>
			<p><a href="%s">%s</a></p>
			<p>С уважением, сервис Производственная гимнастика.</p>
			<p><i>Если вы не запрашивали подтверждение, просто проигнорируйте это сообщение.</i></p>
		</body>
	</html>
	`, userEmail, verifyURL, verifyURL, verifyURL)

	m.AddAlternative("text/html", htmlBody)

	d := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.EmailPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email.
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
