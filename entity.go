package main

type DataMutualGuilds struct {
	ID   string `json:"id"`
	Nick string `json:"nick"`
}

// Структура профиля из ответа
type ProfileData struct {
	MutualGuilds []DataMutualGuilds `json:"mutual_guilds"`
}
