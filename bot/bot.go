package bot

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"task/database"
	"task/model"
	"time"
)

func StartBot(key string) {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		log.Fatalf("failed to create new bot api: %s\n", err)
	}
	bot.Debug = false

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatalf("failed to get update chan: %s\n", err)
	}
	log.Printf("bot started with authorized user: %s\n", bot.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			command := update.Message.Command()
			switch command {

			case "start":
				reply := "Welcome to the task manager bot! You can manage your tasks using the following commands:\n\n" +
					"/addtask <task>: Add a new task\n" +
					"/listtasks: List all tasks\n" +
					"/removetask <task-id>: Remove a task by ID"
				sendMessage(bot, update.Message.Chat.ID, reply)

			case "addtask":
				task := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/addtask"))
				if task != "" {

					lastTaskNum := 0

					// first we need to get the last number of task
					lastTaskId, err := database.Task.FindLastRecord(update.Message.Chat.ID)
					if err != nil {
						if !errors.Is(err, gorm.ErrRecordNotFound) {
							log.Printf("failed to find last number for %v, chat\n", update.Message.Chat.ID)
							sendMessage(bot, update.Message.Chat.ID, "server error. please try again later")
							continue
						}
					} else {
						// unmarshal it to get last number of tasks
						lastTaskNum = unmarshalTaskId(lastTaskId)
					}

					// marshal it into taskId with one added into last task number
					currentTaskId := marshalTaskId(update.Message.Chat.ID, lastTaskNum+1)

					taskRecord := model.Task{
						Id:          currentTaskId,
						CreatedAt:   time.Now(),
						ChatId:      update.Message.Chat.ID,
						Description: task,
					}

					if err := database.Task.CreateRecord(&taskRecord); err != nil {
						log.Printf("failed to task into database, %s\n", err)
						sendMessage(bot, update.Message.Chat.ID, "server error. please try again later")
						continue
					}

					reply := "Task added successfully: " + task
					sendMessage(bot, update.Message.Chat.ID, reply)
				} else {
					reply := "Please provide a task description."
					sendMessage(bot, update.Message.Chat.ID, reply)
				}
			case "listtasks":
				tasks, err := database.Task.FindRecords(update.Message.Chat.ID)
				if err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						log.Printf("failed to find records: %s\n", err)
						sendMessage(bot, update.Message.Chat.ID, "server error. please try again later")
						continue
					}
					reply := "You haven't added any task yet!!"
					sendMessage(bot, update.Message.Chat.ID, reply)
					continue
				}

				reply := "Here are your tasks:\n"
				secondStr := fmt.Sprintf("%-5s%-8s%s\n", "No.", "TaskID", "Task")
				reply = reply + secondStr
				for i, t := range tasks {
					taskNum := unmarshalTaskId(t.Id)
					str := fmt.Sprintf("%-5d%-8s%s\n", i+1, fmt.Sprintf("%04d", taskNum), t.Description)
					reply = reply + str
				}

				sendMessage(bot, update.Message.Chat.ID, reply)
			case "removetask":
				// Handle the /removetask command
				taskIDStr := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/removetask"))
				if taskIDStr != "" {
					taskNum, err := strconv.Atoi(taskIDStr)
					if err != nil {
						sendMessage(bot, update.Message.Chat.ID, "Invalid task ID.")
						continue
					}
					taskId := marshalTaskId(update.Message.Chat.ID, taskNum)
					if err := database.Task.DeleteRecords(taskId); err != nil {
						log.Printf("failed to delete task, %s\n", err)
						sendMessage(bot, update.Message.Chat.ID, "server error. please try again later")
						return
					}
					reply := "Task removed successfully: " + strconv.Itoa(taskNum)
					sendMessage(bot, update.Message.Chat.ID, reply)
				} else {
					sendMessage(bot, update.Message.Chat.ID, "Please provide a task ID.")
				}
			default:
				sendMessage(bot, update.Message.Chat.ID, "Invalid command.")
			}
		}
	}
	return
}

func sendMessage(bot *tgbotapi.BotAPI, id int64, reply string) {
	msg := tgbotapi.NewMessage(id, reply)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

// marshal chatId and taskNumbet into taskId (chatId:taskNum)
func marshalTaskId(chatId int64, taskNum int) string {
	str := fmt.Sprintf("%d:%04d", chatId, taskNum)
	return str
}

// unmarshal taskID(chatId:taskNum) and return taskNum
func unmarshalTaskId(taskId string) int {
	str := strings.Split(taskId, ":")
	num, _ := strconv.Atoi(str[1])
	return num
}
