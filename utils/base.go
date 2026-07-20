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
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
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

func (b *BaseUtil) SafeGet(row map[string]any, key string, defaultVal any) any {
	if row == nil {
		return defaultVal
	}
	if val, ok := row[key]; ok && val != nil {
		return val
	}
	return defaultVal
}

func (b *BaseUtil) GetQueryParam(query url.Values, key string, defaultVal string) string {
	val := query.Get(key)
	if len(val) > 0 {
		return val
	}
	return defaultVal
}

func (b *BaseUtil) Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call: must have even number of arguments")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func (b *BaseUtil) Merge(base interface{}, values ...interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	rv := reflect.ValueOf(base)
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			if field.PkgPath != "" {
				continue // unexported
			}
			out[field.Name] = rv.Field(i).Interface()
		}
	} else if m, ok := base.(map[string]interface{}); ok {
		for k, v := range m {
			out[k] = v
		}
	} else {
		out["Data"] = base
	}

	for i := 0; i < len(values); i += 2 {
		k := values[i].(string)
		out[k] = values[i+1]
	}

	return out, nil
}

func (b *BaseUtil) Coalesce(vals ...interface{}) interface{} {
	for _, v := range vals {
		if b.isTruthy(v) {
			return v
		}
	}
	return nil
}

func (b *BaseUtil) isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case string:
		return val != ""
	case bool:
		return val
	case int:
		return val != 0
	case int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0
	case []interface{}:
		return len(val) > 0
	case []string:
		return len(val) > 0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return rv.Len() > 0
		case reflect.Ptr, reflect.Interface:
			return !rv.IsNil()
		}
	}

	return true
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
