package H

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	mrand "math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	Tinify "github.com/gwpp/tinify-go/tinify"
	jsoniter "github.com/json-iterator/go"
	"github.com/k0kubun/pp"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var P = fmt.Println
var Format = fmt.Sprintf
var Pretty = pp.Println

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
	return fmt.Sprintf("%f", f)

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

// Tif is a simple implementation of the dear ternary IF operator
func Tif(condition bool, tifThen, tifElse interface{}) interface{} {
	if condition {
		return tifThen
	}

	return tifElse
}

func TifStr(condition bool, tifThen, tifElse interface{}) string {
	return Tif(condition, tifThen, tifElse).(string)
}

// MatchesAny returns true if any of the given items matches ( equals ) the subject ( search parameter )
func MatchesAny(search interface{}, items ...interface{}) bool {
	for _, v := range items {
		if fmt.Sprintf("%T", search) == fmt.Sprintf("%T", v) {
			if search == v {
				return true
			}
		}
	}

	return false
}

// StringReplaceAll keeps replacing until there's no more ocurrences to replace.
func StringReplaceAll(original string, replacementPairs ...string) string {
	if original == "" {
		return original
	}

	r := strings.NewReplacer(replacementPairs...)

	for {
		result := r.Replace(original)

		if original != result {
			original = result
		} else {
			break
		}
	}

	return original
}

func Line() {
	fmt.Println("------------------------------------------------------------------")
}

func PL(i ...interface{}) {
	Line()
	P(i...)
	Line()
}

func PBL(i ...interface{}) {

	P(i...)
	Line()
}

func RemoveFromString(original string, removePairs ...string) string {
	var ar []string
	for _, value := range removePairs {
		ar = append(ar, value, "")
	}

	return StringReplaceAll(original, ar...)

}

func UrlString(s string) string {
	s = StringReplaceAll(s, "-", " ")
	return StringReplaceAll(CleanIndexText(s), " ", "-")
}

func UrlStringForFile(s string) string {
	ext := filepath.Ext(s)
	filename := RemoveFromString(s, ext)
	filename = UrlString(filename) + ext
	return filename
}

func CleanIndexText(text string) string {
	text = strings.ToLowerSpecial(unicode.TurkishCase, text)

	text = StringReplaceAll(text, "ı", "i")

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	text, _, _ = transform.String(t, text)

	reg, err := regexp.Compile(`[^a-zA-Z0-9 ]`)
	if err != nil {
		log.Fatal(err)
	}
	text = reg.ReplaceAllString(text, "")

	return text
}

func DownloadFile(url, filepath string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func GetTextBetween(text, strStart, strEnd string) string {
	i1 := strings.Index(text, strStart) + len(strStart)
	i2 := len(text)

	if strEnd != "" {
		i2 = strings.Index(text[i1:], strEnd) + i1
	}

	//error
	if i2 < i1 {
		PL("GetTextBetween ERROR!", "strStart:", strStart, "strEnd:", strEnd)
	}

	return text[i1:i2]
}

func ReadFile(path string) string {
	r, err := ioutil.ReadFile(path)
	if err != nil {
		P(err)
	}
	return string(r)
}

func WriteFile(path, text string) error {

	return ioutil.WriteFile(path, []byte(text), 0644)

}

func AlignText(w int, s ...string) string {
	//return fmt.Sprintf(fmt.Sprintf("%%%ds", w),s)

	specialChars := []string{"₺", "₿"}

	var x string
	var oldLen int

	for i, value := range s {
		space := w - oldLen

		if i == 0 {
			space = 0
		}

		x += Space(space) + value
		oldLen = len(value)

		for _, char := range specialChars {
			if strings.Contains(value, char) {
				oldLen -= 2
			}
		}

	}

	return x
}

func Space(n int) string {
	return strings.Repeat(" ", n)
}

func AppleScriptCommand(s string) string {
	command := exec.Command("/usr/bin/osascript", "-e", s)

	stdout, _ := command.CombinedOutput()

	return string(stdout)

}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsMacos() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func Command(command string) string {

	var cmd *exec.Cmd

	if IsWindows() {
		cmd = exec.Command("cmd", "/C", command)

	} else if IsMacos() {
		cmd = exec.Command("/bin/zsh", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	stdout, _ := cmd.CombinedOutput()

	return string(stdout)
}

func Test() {
	P("test")
}

func FM(format string, a ...interface{}) string {

	//x := FM("num:$num str:$s b:$b",num,s,b)
	//var words []string

	s := Mstr(format)

	r, _ := regexp.Compile("{([^}]*)}")

	allString := r.FindAllString(format, -1)

	for i, value := range allString {
		s.Replace(value, Str(a[i]))
	}

	return s.String()
}

func TitleTurkish(s string) string {
	s = strings.ToLowerSpecial(unicode.TurkishCase, s)

	var arr []string

	for _, value := range strings.Split(s, " ") {
		if len(value) == 0 {
			continue
		}

		arr = append(arr, strings.ToUpperSpecial(unicode.TurkishCase, string([]rune(value)[0]))+string([]rune(value)[1:]))

	}

	return strings.Join(arr, " ")

}

func ReverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func GetRequest(url string, headers ...string) ([]byte, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	for i := 0; i < len(headers)-1; i++ {
		req.Header.Set(headers[i], headers[i+1])
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func PostRequest(postUrl string, postData map[string]string, headers ...string) ([]byte, error) {

	params := url.Values{}
	for key, value := range postData {
		params.Set(key, value)
	}

	pData := strings.NewReader(params.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", postUrl, pData)

	for i := 0; i < len(headers)-1; i++ {
		req.Header.Set(headers[i], headers[i+1])
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	return body, err
}

func parseJsonFromUrl(method, url, path string, i interface{}, postData map[string]string, headers ...string) error {

	var bytedata []byte

	if method == "GET" {
		d, err := GetRequest(url, headers...)
		if err != nil {
			return err
		}
		bytedata = d
	}
	if method == "POST" {
		d, err := PostRequest(url, postData, headers...)
		if err != nil {
			return err
		}
		bytedata = d
	}

	paths := strings.Split(path, ".")

	get := jsoniter.Get(bytedata, paths[0])

	if path != "" {
		bytedata = []byte(get.ToString())
	}

	if len(paths) > 1 {

		for i := 1; i < len(paths); i++ {

			get = get.Get(paths[i])
		}

		bytedata = []byte(get.ToString())
	}

	return jsoniter.Unmarshal(bytedata, i)

}
func ParseJsonFromUrlGET(url, path string, i interface{}, headers ...string) error {

	return parseJsonFromUrl("GET", url, path, i, nil, headers...)

}

func ParseJsonFromUrlPOST(url, path string, i interface{}, postData map[string]string, headers ...string) error {
	return parseJsonFromUrl("POST", url, path, i, postData, headers...)

}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err

}

func FileExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		return true

	}

	return false

}

func NewFolder(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func RandomSeed() {
	mrand.Seed(time.Now().UTC().UnixNano())
}

func RandomInt(min int, max int) int {
	return min + mrand.Intn(max-min)
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandomInt(65, 90))
	}
	return string(bytes)
}

func Svg(path string, size int, stroke float64) string {
	return FM(`<svg id="i-desktop" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="{size}" height="{size}" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="{stroke}">
    <path d="{path}" />
</svg>`, size, size, stroke, path)
}

func TurkishDate(date int64) string {

	theTime := ParseUnix(date).Format("02 January 2006 15:04")
	month := strings.Split(theTime, " ")
	switch month[1] {
	case "January":
		str := strings.Replace(theTime, "January", "Ocak", -1)
		return str
	case "February":
		str := strings.Replace(theTime, "February", "Şubat", -1)
		return str
	case "March":
		str := strings.Replace(theTime, "March", "Mart", -1)
		return str
	case "April":
		str := strings.Replace(theTime, "April", "Nisan", -1)
		return str
	case "May":
		str := strings.Replace(theTime, "May", "Mayıs", -1)
		return str
	case "June":
		str := strings.Replace(theTime, "June", "Haziran", -1)
		return str
	case "August":
		str := strings.Replace(theTime, "August", "Ağustos", -1)
		return str
	case "September":
		str := strings.Replace(theTime, "September", "Eylül", -1)
		return str
	case "October":
		str := strings.Replace(theTime, "October", "Ekim", -1)
		return str
	case "November":
		str := strings.Replace(theTime, "November", "Kasım", -1)
		return str
	case "December":
		str := strings.Replace(theTime, "December", "Aralık", -1)
		return str
	default:
		return "Date Parse Error"
	}
}

func GetLines(s string) []string {
	lines := strings.Split(s, "\n")

	var newLines []string

	for _, line := range lines {
		if line == "" {
			continue
		}
		newLines = append(newLines, line)
	}
	return newLines
}

func GetUniqueFileName(fileName, inFolder string) string {

	files, _ := ioutil.ReadDir(inFolder)
	tryName := fileName

	count := 2
	for {

		var fileExist bool

		for _, k := range files {
			if tryName == k.Name() {

				ext := filepath.Ext(fileName)
				nameWithoutExt := RemoveFromString(fileName, ext)
				tryName = nameWithoutExt + Str(count) + ext
				count++

				fileExist = true
			}
		}

		if !fileExist {
			fileName = tryName
			break
		}

	}

	return fileName

}

func GetRandomNumber(max int) int {
	s := mrand.NewSource(time.Now().Unix())
	r := mrand.New(s)
	return r.Intn(max)
}

func TinyPNG(input, output string) error {

	keys := []string{"hzl3QCNty3fZDkl70ThZPqnMSpZ4Yb02", "LRlR8p3Tk2LyZ7Zjyby1BL5nds3WMQSH", "l1NFy7xfLbNkXrFvXXCc1wjZ3zklQm2G", "HJSGWrygrj1rN9YgZWyBf1SpPrz1Rj0x", "qkHnGMYDgStdl4dHzl0j66Y3ZSfKcD7H", "n1WK10Z7CBLs4kv4Sf5ZZTJQCgJRj2wl", "zHDhPf78Hxrj2NwlpWl8Tp5bNsddJ3FS", "sdkkDwyfKqM6z8MDgnqBNTQJGcTKfH68", "5S63QBhljN8lrZhbH13PmhNjRwcM3FfC", "LMpcgLy6ypYZ71GYTDdT56pjRM0q16bs", "pLpWFf01LqpxTBV65vffJR2qkqFp7wmP", "Lvdzn423mc5S9yBYzqDLRVM5vwCrvdJ5", "JSzhF5DYS0CMqwbXtP24RQxjwQ5rz6FM", "SkYJ5GQsppGbbrMt7dfxmx2BvKBdx1gj", "5qvySGSTn4dpxmxpc1YLs4TPXr3V6NNY", "SbvHFjgBlYch39LdbQ2RWGNf0r7Jm76Z", "9rxbwLJg8xY6Tr6DCJVZCHD88P7MZZ9p", "Nv945ymkc7bFQYbbSxd9F2nQsrkzg0rL", "ZlcnjDRSK6VHs7zCyQ516FT0pp2LBL3l", "6Q74Yf1D3QfKvvySvY7mb53t0XXn4LnL", "nsDnbVzD3Kg526G4KY3PBh9ypSgvjjGr", "HzQ7mwJshKfCYJzqwQHzNsC5rPstHdrH"}

	Tinify.SetKey(keys[GetRandomNumber(len(keys))])
	source, err := Tinify.FromFile(input)
	if err != nil {

		return err
	}

	err = source.ToFile(output)
	if err != nil {

		return err
	}
	return nil

}
