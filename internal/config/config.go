package config

import (
	"analytics/pkg/utils"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      App        `json:"app"`
	Database ClickHouse `json:"clickhouse"`
}

type App struct {
	Env     string `json:"env"`
	Secret  string `json:"secret"`
	AppID   uint   `json:"app_id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type ClickHouse struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Port     int    `json:"port"`
}

var config Config

func MustLoad(filename string) *Config {
	configFile := utils.FindDirectoryName("configs") + "\\" + filename

	if err := cleanenv.ReadConfig(configFile, &config); err != nil {
		panic(err)
	}

	return &config
}

//// watchConfig для лайв релоада конфига
//func watchConfig(filename string) {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		fmt.Println("Ошибка создания watcher:", err)
//		return
//	}
//	defer watcher.Close()
//
//	err = watcher.Add(filename)
//	if err != nil {
//		fmt.Println("Ошибка добавления файла в watcher:", err)
//		return
//	}
//
//	fmt.Println("Отслеживание конфигурации запущено...")
//
//	for {
//		select {
//		case event, ok := <-watcher.Events:
//			if !ok {
//				return
//			}
//
//			if event.Op&fsnotify.Write == fsnotify.Write {
//				fmt.Println("Обнаружено изменение конфигурации, перезагрузка...")
//				if err := loadConfig(filename); err != nil {
//					fmt.Println("Ошибка загрузки конфигурации:", err)
//				}
//			}
//
//		case err, ok := <-watcher.Errors:
//			if !ok {
//				return
//			}
//			fmt.Println("Ошибка watcher:", err)
//		}
//	}
//}
