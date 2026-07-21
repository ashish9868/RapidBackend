package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ashish9868/rapidbackend/dto"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mail.v2"
)

const (
	DEFAULT_PORT = 7003
)

type BaseUtil struct{}

func NewBaseUtil() *BaseUtil {
	return &BaseUtil{}
}

func (util *BaseUtil) SafeEnvGet(name string, defaultVal string) string {
	val := strings.TrimSpace(os.Getenv(name))
	if len(val) > 0 {
		return val
	}
	return defaultVal
}

func (uitl *BaseUtil) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	switch {
	case !hasUpper:
		return errors.New("password must contain an uppercase letter")
	case !hasLower:
		return errors.New("password must contain a lowercase letter")
	case !hasDigit:
		return errors.New("password must contain a number")
	case !hasSpecial:
		return errors.New("password must contain a special character")
	}

	return nil
}

func (*BaseUtil) HashPassword(p string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(h)
}
func (*BaseUtil) CheckPassword(hash, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(p)) == nil
}

func (*BaseUtil) GetSha256Hash(str string) string {
	h := sha256.Sum256([]byte(str))
	return hex.EncodeToString(h[:])
}

func (*BaseUtil) GetCWD() string {
	// Gets the absolute path of the running executable
	ex, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	return ex
}

func (b *BaseUtil) GetPathFromRoot(path string) string {
	return filepath.Join(b.GetCWD(), path)
}

func (b *BaseUtil) SafeCreateFile(path string, content string) error {

	_, err := os.Stat(path)

	if err == nil {
		return nil // File exists
	}
	if errors.Is(err, fs.ErrNotExist) {
		pathDir := filepath.Dir(path)
		err = b.SafeCreateFolder(pathDir)
		if err == nil {
			err := os.WriteFile(path, []byte(content), 0644)
			if err == nil {
				return nil
			}
			return err
		}
		return err // File explicitly does not exist
	}

	return nil
}

func (b *BaseUtil) SafeCreateFolder(folder string) error {
	root := b.GetCWD()
	folders := append([]string{root}, strings.Split(folder, "/")...)
	fullPath := path.Join(folders...)

	fmt.Println("Creating path", fullPath)
	stat, err := os.Stat(fullPath)

	createFolder := true
	if err == nil {
		if stat.IsDir() {
			createFolder = false
		}
	}

	if createFolder {
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*BaseUtil) SendEmail(settings dto.SMTPSetting, message dto.SMTPMessage) error {
	m := mail.NewMessage()
	m.SetAddressHeader("From", settings.From, settings.FromName)
	m.SetHeader("To", message.To)
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/plain", message.Plain)
	m.AddAlternative("text/html", message.Html)
	for _, cc := range message.CC {
		m.SetHeader("Cc", cc)
	}
	d := mail.NewDialer(
		settings.SMTPHost,
		settings.SMTPPort,
		settings.SMTPUserName,
		settings.SMTPPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Email Send failed", err.Error())
		return err
	}
	return nil
}

func (b *BaseUtil) GenerateRandomHash() (raw, hash, prefix string) {
	id := xid.New().String()
	bytes := make([]byte, 16) // 256-bit
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	raw = "tx_" + base64.RawURLEncoding.EncodeToString([]byte(id)) + "_" + base64.RawURLEncoding.EncodeToString(bytes)

	hash = b.GetSha256Hash(raw)

	prefix = raw[:10]

	return
}

func (b *BaseUtil) WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (b *BaseUtil) Debug(data interface{}) string {
	fmt.Printf("DEBUG: %#v\n", data) // prints to console
	return ""                        // return empty so template doesn't render anything
}

func (b *BaseUtil) ToJSON(data any) *string {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("unable to transform to JSON", err.Error())
		return nil
	}
	str := string(bytes)
	return &str
}

func (b *BaseUtil) FileExists(embed fs.FS, path string) bool {
	_, err := fs.Stat(embed, path)
	if err == nil {
		return true
	}
	return false
}

func (b *BaseUtil) PrintFiles(embed fs.FS) error {
	if gin.Mode() == gin.ReleaseMode {
		return nil
	}
	// Walk the root directory "." to visit every file and folder
	err := fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip printing the root folder itself
		if path == "." {
			return nil
		}

		// Distinguish between files and directories
		if d.IsDir() {
			fmt.Printf("[DIR]  %s\n", path)
		} else {
			fmt.Printf("[FILE] %s\n", path)
		}
		return nil
	})

	return err
}

func (b *BaseUtil) SubFs(embed fs.FS, path string) *fs.FS {
	subFs, err := fs.Sub(embed, path)
	if err == nil {
		return &subFs
	}
	println("Error reading sub fs :", err.Error())
	return nil
}
