package server

import (
	"github.com/fshmidt/rassilki/pkg/service"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type RassilkaListener struct {
	sync.RWMutex
	rasMap          map[int]*RassilkaStatus
	messagesService service.Messages
}

func NewRassilkaListener(messagesService service.Messages) *RassilkaListener {
	return &RassilkaListener{
		RWMutex:         sync.RWMutex{},
		rasMap:          make(map[int]*RassilkaStatus),
		messagesService: messagesService}
}

type RassilkaStatus struct {
	Active    bool
	Started   bool
	Updated   bool
	Recreated bool
}

func (r *RassilkaListener) UpdateStatus(actIds, allIds, updatedIds, recreatedIds []int) {
	r.Lock()
	defer r.Unlock()
	logrus.Println("recrecreated:", recreatedIds, "updated:", updatedIds, "accIds", actIds)
	// Закрытие горутин для удаленных и пересозданных рассылок
	for k, _ := range r.rasMap {
		if r.isDeleted(k, allIds) {
			delete(r.rasMap, k)
		}
	}

	// Остановка горутин для измененных рассылок
	for k, status := range r.rasMap {
		if r.isStopped(k, actIds) {
			status.Active = false
		}
	}
	// Обновление статусов для активных рассылок
	for _, id := range actIds {
		r.activateStatus(id)
	}
	// Обновление статусов для заново сформированных рассылок
	for _, id := range recreatedIds {
		r.changeRecreateStatus(id)
	}
	// Обновление статусов для дополненных рассылок
	for _, id := range updatedIds {
		r.changeUpdateStatus(id)
	}

	// Добавление новых рассылок

	for _, id := range actIds {
		if _, ok := r.rasMap[id]; !ok {
			r.rasMap[id] = &RassilkaStatus{
				Active: true,
			}
			started, err := r.messagesService.IsStarted(id)
			if err != nil {
				logrus.Println(err.Error())
			}
			if started {
				r.rasMap[id].Started = true
			}
			go r.startRassilkaWorker(id)
		}
	}
}

func (r *RassilkaListener) isDeleted(id int, allIds []int) bool {

	for _, val := range allIds {
		if id == val {
			return false
		}
	}
	return true
}

func (r *RassilkaListener) isStopped(id int, actIds []int) bool {
	for _, val := range actIds {
		if id == val {
			return false
		}
	}
	return true
}

func (r *RassilkaListener) activateStatus(id int) {
	if status, ok := r.rasMap[id]; ok {
		status.Active = true
		if isStated, err := r.messagesService.IsStarted(id); isStated && err == nil {
			status.Started = true
		} else if err != nil {
			logrus.Println(err)
		}
	}
}

func (r *RassilkaListener) changeUpdateStatus(id int) {
	if status, ok := r.rasMap[id]; ok {
		status.Updated = true
	}
}

func (r *RassilkaListener) changeRecreateStatus(id int) {
	if status, ok := r.rasMap[id]; ok {
		status.Recreated = true
	}
}

func (r *RassilkaListener) startRassilkaWorker(id int) {

	for {
		r.RLock()

		clients, err := r.messagesService.GetClientsList(id)
		if err != nil {
			logrus.Println("can't get clients list.", err)
		}

		if status, ok := r.rasMap[id]; ok {
			if status.Active && !status.Started {
				logrus.Println("Worker N", id, "is getting clients ids")

				if err := r.messagesService.CreateTable(id, clients); err != nil {
					logrus.Println(err)
				}

				logrus.Println("Worker N", id, "created its message table")

			} else if status.Active && status.Started {

				messagesLeft := 0

				for _, client := range clients {

					message, err := r.messagesService.GetMessage(client.Id, id)
					if err != nil {
						logrus.Println("can't get message struct for client", client.Id, err)
					}

					msgText, err := r.messagesService.GetMessageText(id)
					if err != nil {
						logrus.Println("can't get message text for client")
					}
					clientPhone, err := strconv.Atoi(client.Code + client.Phone)
					if err != nil {
						logrus.Println("can't convert clinetPhone to int", err)
					}
					if message.Status == true {
						continue
					}
					messagesLeft += 1

					err = r.messagesService.SendMessageToSubscriber(message.Id, clientPhone, id, msgText)
					if err != nil {
						logrus.Println("problem with sending message to subscriber", clientPhone, ". error:", err)
					} else {
						logrus.Println("Rassilka N", id, "sending '", msgText, "' to", client.Id, client.Phone)

					}

				}

				if messagesLeft == 0 {
					logrus.Println("Rassilka N", id, "is still active but it has already finished sending all messages. You can add more clients or change/delete rassilka.")
				}
			} else {

				logrus.Println("Rassilka N", id, "is anactive and waits for the beggining.")
			}
		} else {
			if err := r.messagesService.DropTable(id); err != nil {
				logrus.Println(err)
			}
			logrus.Println("Worker N", id, "удалил таблицу и закрыл горутину")
			r.RUnlock()
			return
		}
		r.RUnlock()
		time.Sleep(6 * time.Second)
	}
}

func (r *RassilkaListener) ResetUpdateRecreateStatus(updatedIds, recreatedIds []int) {
	r.Lock()
	defer r.Unlock()
	for _, id := range updatedIds {
		if status, ok := r.rasMap[id]; ok {
			status.Updated = false
		}
	}
	for _, id := range recreatedIds {
		if status, ok := r.rasMap[id]; ok {
			status.Recreated = false
		}
	}
}
