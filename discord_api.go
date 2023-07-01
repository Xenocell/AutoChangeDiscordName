package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// мой идентификатор по умолчанию в строке
const MyStrUID = "@me"

var (
	//URL для получения профиля по идентификатору пользователя
	EndpointProfile = func(uID string) string {
		return fmt.Sprintf("https://discord.com/api/v9/users/%s/profile", uID)
	}
	//URL для изменения ника на сервере
	EndpointChangeName = func(gID, uID string) string {
		return fmt.Sprintf("https://discord.com/api/v9/guilds/%s/members/%s", gID, uID)
	}
)

type apiClient struct {
	authToken string
	myUID     string
	client    *http.Client
}

// инициализация структуры
func NewApiClient(authToken string, myUID string) *apiClient {
	return &apiClient{
		authToken: authToken,
		client:    &http.Client{},
	}
}

// Запрос на изменение ника по иду гильдии
func (ac *apiClient) ChangeMyNickOnTheGuild(gID, nick string) error {

	payload, err := json.Marshal(map[string]interface{}{
		"nick": nick,
	})
	if err != nil {
		return fmt.Errorf("main - ChangeMyNickOnTheGuild - json.Marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, EndpointChangeName(gID, MyStrUID), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("main - ChangeMyNickOnTheGuild - http.NewRequest: %w", err)
	}

	req.Header.Add("Authorization", ac.authToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ac.client.Do(req)
	if err != nil {
		return fmt.Errorf("main - ChangeMyNickOnTheGuild - ac.client.Do: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("main - ChangeMyNickOnTheGuild: unexpected error (status_code - %d)", resp.StatusCode)
	}

	return nil
}

// запрос на получения текущего профиля
func (ac *apiClient) GetMyProfile() (ProfileData, error) {

	var profileData ProfileData

	req, err := http.NewRequest(http.MethodGet, EndpointProfile(ac.myUID), nil)
	if err != nil {
		return profileData, fmt.Errorf("main - GetMyProfile - http.NewRequest: %w", err)
	}

	req.Header.Add("Authorization", ac.authToken)

	resp, err := ac.client.Do(req)
	if err != nil {
		return profileData, fmt.Errorf("main - GetMyProfile - ac.client.Do: %w", err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&profileData)
	if err != nil {
		return profileData, fmt.Errorf("main - GetMyProfile - json.NewDecoder: %w", err)
	}

	return profileData, nil
}
