package H

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/savsgio/atreugo/v11"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"time"
)
import "golang.org/x/crypto/bcrypt"

var P = log.Println
var Format = fmt.Sprintf

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func F(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func FtoStr(f float64) string {
	return fmt.Sprintf("%v", f)

}

func Int(i interface{}) int {

	if i == nil {
		return 0
	}
	result, err := strconv.Atoi(i.(string))
	if err != nil {
		return 0
	}
	return result
}

func IntToStr(i int) string {
	return strconv.Itoa(i)
}

func PrintMemUsage() {

	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func Now() int64 {
	return time.Now().Unix()
}

func ParseUnix(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func TimeFormat(date int64) string {
	return ParseUnix(date).Format("2 Jan 2006 15:04:05")
}

func Profit(val1 float64, val2 float64) float64 {

	p := 100 * (val2 - val1) / val1
	return Fixed(p, 2)
}

func Fixed(f float64, prec int) float64 {
	x := math.Pow10(prec)
	return math.Round(f*x) / x
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func Truncate(text string, length int) string {

	if len(text) < length {
		return text
	}
	return text[0:length] + "..."
}

func Str(i interface{}) string {

	return fmt.Sprintf("%v", i)
}

// ============================================== FORM HELPER ======================================================
type Form struct {
	Ctx *atreugo.RequestCtx
}

func (f Form) Get(key string) string {
	return string(f.Ctx.FormValue(key))
}

func (f Form) GetFloat(key string) float64 {
	return F(f.Get(key))
}

func (f Form) GetInt(key string) int {
	return Int(f.Get(key))
}

func (f Form) PrintPostData() {
	P(string(f.Ctx.PostBody()))
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s | %s", name, elapsed)
}

var errAesEncrypt = errors.New("Aes decrypt error")

func AESEncrypt(key, text string) (string, error) {

	c, err := aes.NewCipher([]byte(key))
	// if there are any errors, handle them
	if err != nil {
		return "", err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return "", errAesEncrypt

	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errAesEncrypt

	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return string(gcm.Seal(nonce, nonce, []byte(text), nil)), err

}

const AesKey = "82!@#$F^&(_+(*^cd^~Z?a$%^&sVxT)*"

var errAesDecrypt = errors.New("Aes decrypt error")

func AESDecrypt(key, text string) (string, error) {

	if text == "" {
		return "", errors.New("empty string")
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errAesDecrypt
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", errAesDecrypt

	}

	t := []byte(text)

	nonceSize := gcm.NonceSize()
	if len(t) < nonceSize {
		return "", errAesDecrypt
	}

	nonce, ciphertext := t[:nonceSize], t[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errAesDecrypt
	}

	return string(plaintext), err

}

func Decrypt(s string) (string, error) {

	base64Dec, err := base64.StdEncoding.DecodeString(s)

	if err != nil {
		return "", err
	}
	return AESDecrypt(AesKey, string(base64Dec))

}

func Encrypt(s string) (string, error) {
	s, err := AESEncrypt(AesKey, s)

	base64Str := base64.StdEncoding.EncodeToString([]byte(s))

	return base64Str, err
}

func GetRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	return contents, nil
}
