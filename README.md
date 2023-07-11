# Telegram-taskmanager-bot

A telegram task manager bot written in golang.

This project used `"github.com/go-telegram-bot-api/telegram-bot-api"`
to communicate with telegram. This Project used postgres DB to store datas.
This Project does not store any kind of personal data (except for the tasks and chatID 
mainly used for retrieval).You can enable debugging in `bot/bot.go`

You can set variable either by ENV or `config.json`. check example_config.json.

ENV LISTS
```
BOT_API_KEY:"botkey"
DATABASE_HOST:"dbhost
DATABASE_USER:"dbuser"
DATABASE_:PASSWORD:"dbpassword
DATABASE_NAME:"dbname"
DATABASE_PORT: 1234
DATABASE_SSLMODE: disable
```

Project Dependencies
1. `gorm "gorm.io/gorm" (for tasks storage)`
2. `viper "github.com/spf13/viper"(for getting envs)`
3. `validator "github.com/go-playground/validator/v10" (for validating config)`
4. `telegrambotapi  "github.com/go-telegram-bot-api/telegram-bot-api"`

** Note - some of the dependencies like viper and validator can be replace with smaller dependencies