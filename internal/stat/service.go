package stat

import (
	"go/adv-demo/pkg/event"
	"log"
)

type StatServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

type StatService struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			id, ok := msg.Payload.(uint)
			if !ok {
				log.Println("Bad data in " + event.EventLinkVisited + " " + msg.Payload.(string))
				continue
			}
			s.StatRepository.AddClick(id)
		}
	}
}
