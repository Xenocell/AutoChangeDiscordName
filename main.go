package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Инициализируем конфиг
	cfg, err := InitConfig()
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	//Создаем структуру ядра
	kernel, err := InitialKernel(cfg.AuthToken, cfg.MyUID, cfg.MonitGuilds)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}
	//запускаем поток мониторинга смены ника
	kernel.StartMonitChangeNick()
	//создаем канал
	done := make(chan os.Signal, 1)
	//подписываемся на уведомление о прерывании ctrl + c
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	//ожидаем сигнал в канале
	<-done
	fmt.Println("Program termination detected!")
	//Останавливаем поток мониторинга смены ника
	kernel.StopMonitChangeNick()
}
