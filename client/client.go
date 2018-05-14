package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PClient struct {
	Settings     Settings
	saveCallback func() error
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Settings struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    *Token `json:"token"`
	UserId   uint64 `json:"user_id"`
}

func (c *PClient) GetUserId() uint64 {
	if uid := c.Settings.UserId; uid > 0 {
		return uid
	}

	// c.Index() ...
	c.Settings.UserId = 1663085

	c.saveCallback()
	return c.Settings.UserId
}

func NewClient(set Settings, saveCallback func() error) (*PClient, error) {
	c := PClient{Settings: set, saveCallback: saveCallback}

	if c.Settings.Token == nil {
		err := c.signIn()
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (c *PClient) GetPage(page string) (string, error) {
	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		return "", err
	}

	if c.Settings.Token.TokenType != "bearer" {
		err = c.signIn()
		if err != nil {
			return "", err
		}
		if c.Settings.Token.TokenType != "bearer" {
			return "", fmt.Errorf("token type %s !=  bearer", c.Settings.Token.TokenType)
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Settings.Token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	ans, _ := ioutil.ReadAll(resp.Body)

	return string(ans), nil
}

func (c *PClient) signIn() error {
	body := url.Values{}
	body.Add("grant_type", "password")
	body.Add("username", c.Settings.Username)
	body.Add("password", c.Settings.Password)
	reqBody := body.Encode()

	reqUrl := "https://na.preva.com/exerciser-api//oauth/token"

	req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	req.SetBasicAuth("precorandroidprod", "TuJveE9LEc1NicusEw1Y")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ans := Token{}
	json.Unmarshal(answer, &ans)

	c.Settings.Token = &ans
	err = c.saveCallback()
	if err != nil {
		return err
	}
	return nil
}
