package bootstrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

// Log in as the user that wants to install this agent
func Login(cfg *Config) (string, string, error) {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	if email == "" {
		return "", "", errors.New("Need user's email to login, to register agent into the system")
	}
	if password == "" {
		return "", "", errors.New("Need user's password to login, to register agent into the system")
	}
	body := map[string]string{
		"email":    email,
		"password": password,
	}

	b, _ := json.Marshal(body)

	resp, err := http.Post(cfg.APIURL+"/auth/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var res struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", "", err
	}

	return res.AccessToken, res.RefreshToken, nil
}

func Logout(apiURL, refreshToken string) error {
	body := map[string]string{
		"refresh_token": refreshToken,
	}

	b, _ := json.Marshal(body)

	_, err := http.Post(apiURL+"/auth/logout", "application/json", bytes.NewBuffer(b))
	return err
}
