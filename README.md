# Telegram-taskmanager-bot

A telegram task manager bot written in golang.

This project used `"github.com/go-telegram-bot-api/telegram-bot-api"`
to communicate with telegram. This Project used postgres DB to store datas.
This Project does not store any kind of personal data (except for the tasks and chatID 
mainly used for retrieval)

Project Dependencies
1. `gorm "gorm.io/gorm" (for tasks storage)`
2. `viper "github.com/spf13/viper"(for getting envs)`
3. `validator "github.com/go-playground/validator/v10" (for validating config)`
4. `telegrambotapi  "github.com/go-telegram-bot-api/telegram-bot-api"`

** Note - some of the dependencies like viper and validator can be replace with smaller dependencies