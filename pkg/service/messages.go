package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fshmidt/rassilki"
	"github.com/fshmidt/rassilki/pkg/repository"
	"net/http"
)

type MessagesService struct {
	repo  repository.Messages
	Token string
}

func NewMessagesService(repo repository.Messages, jwt string) *MessagesService {
	return &MessagesService{
		repo:  repo,
		Token: jwt,
	}
}

func (s *MessagesService) CreateTable(id int, clients []rassilki.Client) error {
	return s.repo.CreateTable(id, clients)
}

func (s *MessagesService) DropTable(id int) error {
	return s.repo.DropTable(id)
}

func (s *MessagesService) GetClientsList(rassilkaId int) ([]rassilki.Client, error) {
	return s.repo.GetClientsList(rassilkaId)
}

func (s *MessagesService) GetMessageText(rassilkaId int) (string, error) {
	return s.repo.GetMessageText(rassilkaId)
}

func (s *MessagesService) GetMessage(clientId, rassilkaId int) (rassilki.Message, error) {
	return s.repo.GetMessage(clientId, rassilkaId)
}

func (s *MessagesService) IsStarted(id int) (bool, error) {
	return s.repo.IsStarted(id)
}

func (s *MessagesService) RenewTable(input rassilki.UpdateRassilka, id int) error {
	return s.repo.RenewTable(input, id)
}

func (s *MessagesService) SendMessageToSubscriber(msgID, clientPhone, rassilkaId int, text string) error {
	// Подготовка данных для запроса
	requestBody := map[string]interface{}{
		"id":    msgID,
		"phone": clientPhone,
		"text":  text,
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Создание запроса
	req, err := http.NewRequest("POST", fmt.Sprintf("https://probe.fbrq.cloud/v1/send/%d", msgID), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Добавление заголовка Authorization с использованием JWT токена
	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Обработка ответа
	if resp.StatusCode != http.StatusOK {
		// Обработка ошибки в случае неуспешного запроса
		return fmt.Errorf("неудачный запрос: %d", resp.StatusCode)
	} else {
		s.repo.UpdateMessageStatus(msgID, rassilkaId)
	}
	return nil
}

func (s *MessagesService) GetRassilkiReview(ids []int) ([]rassilki.RassilkaReview, error) {

	return s.repo.GetRassilkiReview(ids)
}
func (s *MessagesService) GetRassilkaReviewById(id int) (rassilki.RassilkaReview, error) {

	return s.repo.GetRassilkaReviewById(id)
}
