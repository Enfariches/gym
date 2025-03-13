package utils

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jordan-wright/email"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
)

func Message(status string, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, src)
	return err
}

func SendEmailCode(addressTo string) (uint, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	addressFrom := os.Getenv("EMAIL_ADDRESS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	password := os.Getenv("EMAIL_PASS")

	rand.Seed(time.Now().UnixNano())
	code := uint(rand.Intn(999_999-111_111+1) + 111_111)
	message := []byte(fmt.Sprintf("Код подтверждения для входа - %d", code))

	e := email.NewEmail()
	e.From = fmt.Sprintf("Гимнастика <%s>", addressFrom)
	e.To = []string{addressTo}
	e.Subject = "Код подтверждения"
	e.Text = message

	err = e.Send(smtpHost+":"+smtpPort, smtp.PlainAuth("", addressFrom, password, smtpHost))
	if err != nil {
		return 0, err
	}
	return code, nil
}
