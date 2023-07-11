package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"task/config"
	"task/model"
)

var db *gorm.DB

func ConnectDatabase(c *config.Config) {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host, c.User, c.Password,
		c.Name, c.Port, c.Sslmode)
	d, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database, error: %s", err)
	}
	if err := d.AutoMigrate(&model.Task{}); err != nil {
		log.Fatalf("failed to migrate tables: %s\n", err)
	}
	db = d
	log.Println("connected to database successfully")
}

type TaskPersistenceInterface interface {
	CreateRecord(task *model.Task) error
	FindRecords(chatId int64) ([]model.Task, error)
	DeleteRecords(taskId string) error
	FindLastRecord(chatId int64) (string, error)
}

var Task TaskPersistenceInterface

type TaskObj struct{}

func (TaskObj) CreateRecord(task *model.Task) error {
	result := db.Create(task)
	return result.Error
}

func (TaskObj) FindRecords(chatId int64) ([]model.Task, error) {
	var tasks []model.Task
	result := db.Where("chat_id = ?", chatId).Find(&tasks)
	if result.RowsAffected == 0 {
		return tasks, gorm.ErrRecordNotFound
	}
	return tasks, result.Error
}

func (TaskObj) DeleteRecords(taskId string) error {
	result := db.Where("id = ?", taskId).Delete(&model.Task{})
	return result.Error
}

func (TaskObj) FindLastRecord(chatId int64) (string, error) {
	var task model.Task
	result := db.Order("created_at desc").Where("chat_id = ?", chatId).First(&task)
	return task.Id, result.Error
}

func init() {
	Task = TaskObj{}
}
