package services

import (
	"messaging-service/database"
	"messaging-service/dto"
	"messaging-service/models"

	"github.com/sirupsen/logrus"
)

func CreateMessage(senderID, receiverID uint, content string) error {
	message := models.Message{SenderID: senderID, ReceiverID: receiverID, Content: content}
	return database.DB.Create(&message).Error
}

func AcceptMessage(messageID uint) error {
	return database.DB.Model(&models.Message{}).Where("id = ?", messageID).Update("is_accepted", true).Error
}

// GetMessageHistoryForUser retrieves the message history between two users
func GetMessageHistoryForUser(senderUserName string, receiverUserName string) ([]dto.MessagesResponse, error) {
	var senderID uint
	if err := database.DB.Model(&models.User{}).
		Where("username = ?", senderUserName).
		Pluck("id", &senderID).Error; err != nil {
		logrus.Error("error getting sender", err)
		return nil, err
	}

	var receiverID uint
	if err := database.DB.Model(&models.User{}).
		Where("username = ?", receiverUserName).
		Pluck("id", &receiverID).Error; err != nil {
		logrus.Error("error getting receiver", err)
		return nil, err
	}

	var messages []models.Message
	if err := database.DB.
		Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		logrus.Error("error getting messages", err)
		return nil, err
	}

	// Create a slice to hold the DTOs
	var responseMessages []dto.MessagesResponse

	// Map each message to the DTO
	for _, msg := range messages {
		var senderName, receiverName string

		// Fetch sender name
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", msg.SenderID).
			Pluck("username", &senderName).Error; err != nil {
			logrus.Error("error getting sender username", err)
			return nil, err
		}

		// Fetch receiver name
		if err := database.DB.Model(&models.User{}).
			Where("id = ?", msg.ReceiverID).
			Pluck("username", &receiverName).Error; err != nil {
			logrus.Error("error getting receiver username", err)
			return nil, err
		}

		// Append the DTO to the response slice
		responseMessages = append(responseMessages, dto.MessagesResponse{
			ID:         msg.ID,
			CreatedAt:  msg.CreatedAt,
			UpdatedAt:  msg.UpdatedAt,
			Sender:     senderName,
			Receiver:   receiverName,
			Content:    msg.Content,
			IsAccepted: msg.IsAccepted,
		})
	}

	return responseMessages, nil
}
