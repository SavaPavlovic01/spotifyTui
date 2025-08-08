package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// rand.Int return error if rand.Read returns one.
// rand.Read says that it never returns an error.
// if it returns an error serious os problem.
func getRandomInt(max int) int64 {
	val, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return val.Int64()
}

func getRandomString(length int) []byte {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	str := make([]byte, length)
	for i := range length {
		index := getRandomInt(len(letters))
		str[i] = letters[index]
	}
	return str
}

func hashString(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func toBase64(s []byte) string {
	return base64.RawURLEncoding.EncodeToString(s[:])
}

func getScopes() string {
	scopes := []string{"streaming", "user-modify-playback-state", "user-read-playback-state"}
	var builder strings.Builder
	builder.Grow(len(scopes))
	for _, scope := range scopes {
		builder.WriteString(scope)
		builder.WriteRune(' ')
	}
	res := builder.String()
	return res[:len(res)-1]
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}

func auth(code string) error {
	headers := map[string]string{
		"client_id":             "c70ac4ad86994ae2aebe9b0da5d708eb",
		"response_type":         "code",
		"redirect_uri":          "http://127.0.0.1:54891/callback",
		"state":                 "idk", // TODO: replace with generateString
		"scope":                 getScopes(),
		"code_challenge_method": "S256",
		"code_challenge":        code,
	}

	params := url.Values{}

	for key, val := range headers {
		params.Add(key, val)
	}
	authUrl := "https://accounts.spotify.com/authorize"
	url := authUrl + "?" + params.Encode()
	err := openBrowser(url)
	if err != nil {
		panic(err)
	}
	return nil
}

func makePostRequestForTokens(params map[string]string) FreshToken {
	form := url.Values{}
	for key, val := range params {
		form.Add(key, val)
	}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		errorMsg, _ := io.ReadAll(resp.Body)

		fmt.Println(string(errorMsg))
		panic("ERROR")
	}
	defer resp.Body.Close()
	var Payload struct {
		Access_token  string  `json:"access_token"`
		Token_type    string  `json:"token_type"`
		Scope         string  `json:"scope"`
		Expires_in    float64 `json:"expires_in"`
		Refresh_token string  `json:"refresh_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&Payload)
	if err != nil {
		panic(err)
	}
	expires_at := time.Now().Add(time.Second * time.Duration(int(Payload.Expires_in))).UnixMilli()
	return FreshToken{AccessToken: Payload.Access_token, RefreshToken: Payload.Refresh_token, ExpiresIn: int(expires_at)}
}

func getToken(code string, verifier string) FreshToken {
	headers := map[string]string{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  "http://127.0.0.1:54891/callback",
		"client_id":     "c70ac4ad86994ae2aebe9b0da5d708eb",
		"code_verifier": verifier,
	}
	return makePostRequestForTokens(headers)
}

type FreshToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func startHttpServer(tokenChan chan FreshToken, verifier string) func() error {
	mux := http.NewServeMux()
	httpWaitForToken := func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		tokens := getToken(code, verifier)
		tokenChan <- tokens
		close(tokenChan)
	}
	mux.HandleFunc("/callback", httpWaitForToken)
	server := &http.Server{
		Addr:    "127.0.0.1:54891",
		Handler: mux,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	}
}

func writeTokens(token FreshToken) error {
	data, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}
	path, err := getTokenFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write token to file")
	}
	return nil
}

func getTokenFilePath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Failed to get path")
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "data.tok"), nil
}

func LoadTokens() (*FreshToken, error) {
	path, err := getTokenFilePath()
	fmt.Print(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read tokens, doing auth flow")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read tokens, doing auth flow")
	}
	var tokens FreshToken
	err = json.Unmarshal(data, &tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to read tokens, doing auth flow")
	}
	return &tokens, nil

}

func (ft *FreshToken) RefreshTokens() {
	headers := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": ft.RefreshToken,
		"client_id":     "c70ac4ad86994ae2aebe9b0da5d708eb",
	}
	newTokens := makePostRequestForTokens(headers)
	ft.AccessToken = newTokens.AccessToken
	ft.RefreshToken = newTokens.RefreshToken
	ft.ExpiresIn = newTokens.ExpiresIn
	writeTokens(*ft)
}

func (ft *FreshToken) Expired() bool {
	return true //return int64(ft.ExpiresIn)-time.Now().UnixMilli() < 60_000
}

func InitialAuth() *FreshToken {
	tokenChannel := make(chan FreshToken)
	verifier := getRandomString(128)
	code := toBase64(hashString(verifier))
	shutdownServer := startHttpServer(tokenChannel, string(verifier))
	auth(code)
	var tokens FreshToken
	for token := range tokenChannel {
		tokens = token
	}
	err := shutdownServer()
	if err != nil {
		fmt.Println(err)
	}
	writeTokens(tokens)
	return &tokens
}
