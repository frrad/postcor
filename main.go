package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/frrad/settings"
)

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Settings struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	set := Settings{
		Username: "your@email.com",
		Password: "secret12312",
	}
	settings.NewSettings(&set, []string{"~/.postcor"})

	token, _ := signIn(set.Username, set.Password)
	log.Printf("%+v\n", token)

	getPage("https://na.preva.com/exerciser-api//exerciser-account", token)
}

func getPage(page string, tok token) error {
	req, err := http.NewRequest("GET", page, nil)

	if tok.TokenType != "bearer" {
		return fmt.Errorf("token type %s !=  bearer", tok.TokenType)
	}

	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	ans, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(ans))

	return nil
}

func signIn(username, password string) (token, error) {

	body := url.Values{}
	body.Add("grant_type", "password")
	body.Add("username", username)
	body.Add("password", password)
	reqBody := body.Encode()

	reqUrl := "https://na.preva.com/exerciser-api//oauth/token"

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(reqBody))
	if err != nil {
		return token{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.SetBasicAuth("precorandroidprod", "TuJveE9LEc1NicusEw1Y")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return token{}, err
	}

	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return token{}, err
	}

	ans := token{}
	json.Unmarshal(answer, &ans)
	return ans, nil
}
