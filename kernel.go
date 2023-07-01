package main

import (
	"fmt"
	"time"
)

const delayUpdate time.Duration = 10

type kernel struct {
	apiClient    *apiClient
	currentNicks map[string]string //key - gUID, value - nick
	done         chan struct{}
}

// инициализируем структуру ядра
func InitialKernel(authToken, myUID string, monitGuilds []string) (*kernel, error) {
	//создаем апи клиент
	apiClient := NewApiClient(authToken, myUID)
	//получаем наш текущий профиль
	profileData, err := apiClient.GetMyProfile()
	if err != nil {
		return nil, fmt.Errorf("main - InitialKernel - apiClient.GetMyProfile: %w", err)
	}
	//создаем массив для записи текущих ников
	currentNicks := make(map[string]string)
	//записываем в созданный массив текущие ники в соответсвии с указанными серверами в параметре monitGuilds
	for _, monitGuild := range monitGuilds {
		for _, mutualGuild := range profileData.MutualGuilds {
			if monitGuild == mutualGuild.ID {
				currentNicks[mutualGuild.ID] = mutualGuild.Nick
			}
		}
	}

	return &kernel{
		apiClient:    NewApiClient(authToken, myUID),
		currentNicks: currentNicks,
		done:         make(chan struct{}),
	}, nil
}

// Запуск отдельного потока
func (k *kernel) StartMonitChangeNick() {
	fmt.Println("Started monitoring the change of nickname")
	go func() {
		for { //Бесконечный цикл
			select {
			case <-k.done: //Ожидаем сигнал о завершении потока
				return
			default: //По умолчанию проверяем смену ника
				err := k.checkChangeNick()
				if err != nil {
					fmt.Println(err)
				}
			}
			// делей на 10 секунд
			time.Sleep(delayUpdate * time.Second)
		}
	}()
}

// Завершаем поток
func (k *kernel) StopMonitChangeNick() {
	k.done <- struct{}{} //Посылаем в канал пустую структуру для остановки потока
	fmt.Println("Monitoring of nickname changes is not stopped")
}

func (k *kernel) checkChangeNick() error {
	//Получаем текущий профиль
	profileData, err := k.apiClient.GetMyProfile()
	if err != nil {
		return fmt.Errorf("checkChangeNick - k.apiClient.GetMyProfile: %w", err)
	}

	for guid, nick := range k.currentNicks {
		for _, mutualGuild := range profileData.MutualGuilds {
			//Проверяем на изм. ника
			if guid == mutualGuild.ID && nick != mutualGuild.Nick {
				fmt.Println("Nick change detected on server " + mutualGuild.ID)
				//делаем запрос на изм. ника, если ник был изменен
				err := k.apiClient.ChangeMyNickOnTheGuild(guid, nick)
				if err != nil {
					fmt.Printf("error updating a nickname on the guild: %s (%s)", guid, err.Error())
				} else {
					getNick := func() string {
						if nick == "" {
							return "by default"
						}
						return nick
					}
					fmt.Println("Nick successfully changed to " + getNick())
				}
				//делей 2 секунды
				time.Sleep(2 * time.Second)
			}
		}
	}

	return nil
}
