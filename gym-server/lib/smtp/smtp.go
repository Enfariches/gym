package smtp

import (
	"fmt"
	"health/internal/config"

	"gopkg.in/mail.v2"
)

func SendAuthMail(smtpConfig config.SMTPConfig, authToken, userEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", smtpConfig.EmailSender)
	m.SetHeader("To", userEmail)
	m.SetAddressHeader("Cc", userEmail, "")
	m.SetHeader("Subject", "Gym: Подтвердждения email")

	verifyURL := fmt.Sprintf("http://localhost:8080/VerifyRegister?auth_token=%s", authToken)
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
			<p>Или введите токен в ручную:</p>
			<p>%s</p>
			<p>С уважением, сервис Производственная гимнастика.</p>
			<p><i>Если вы не запрашивали подтверждение, просто проигнорируйте это сообщение.</i></p>
		</body>
	</html>
	`, userEmail, verifyURL, authToken)

	m.SetBody("text/html", htmlBody)

	d := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.EmailPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendResetMail(smtpConfig config.SMTPConfig, resetToken, userEmail string) error {
	m := mail.NewMessage()
	m.SetHeader("From", smtpConfig.EmailSender)
	m.SetHeader("To", userEmail)
	m.SetAddressHeader("Cc", userEmail, "")
	m.SetHeader("Subject", "Gym: Смена пароля пользователя")

	verifyURL := fmt.Sprintf("http://localhost:8080/VerifyChangePassword?reset_token=%s", resetToken)
	htmlBody := fmt.Sprintf(`
	<html>
		<body>
			<p>Уважаемый пользователь!</p>
			<p>Сервис производственная гимнастика запрашивает подтверждение, что email <b>%s</b> действительно требует смену пароля.</p>
			<p>
				<a href="%s" style="
					display: inline-block;
					padding: 12px 24px;
					font-size: 16px;
					color: #ffffff;
					background-color: #007BFF;
					text-decoration: none;
					border-radius: 6px;
				">Подтвердить смену пароля</a>
			</p>
			<p>Или скопируйте ссылку в браузер:</p>
			<p><a href="%s">%s</a></p>
			<p>С уважением, сервис Производственная гимнастика.</p>
			<p><i>Если вы не запрашивали подтверждение, просто проигнорируйте это сообщение.</i></p>
		</body>
	</html>
	`, userEmail, verifyURL, verifyURL, verifyURL)

	m.SetBody("text/html", htmlBody)

	d := mail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.EmailPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
