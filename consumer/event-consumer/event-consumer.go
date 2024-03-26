package event_consumer

import (
	"TelegramGOBot/events"
	"errors"
	"fmt"
	"log"
	"time"
)

type Consumer struct {
	fetcher       events.Fetcher
	processor     events.Processor
	batchSize     int
	backupStorage map[string]events.Event
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	retryInterval := 1
	maxRetries := 3
	retries := 0

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			retries++
			if retries > maxRetries {
				log.Println("Reached maximum retries, moving on.")
				break
			}
			log.Printf("Retrying in %d seconds...", retryInterval)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			retryInterval *= 2
			continue
		}

		if len(gotEvents) == 0 {
			log.Printf("Retrying in %d seconds...", retryInterval)
			time.Sleep(time.Duration(retryInterval) * time.Second)
			retryInterval *= 2
			continue
		}

		retries = 0
		retryInterval = 1

	}
	return nil
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("Got new event: %s", event.Text)
		retries := 0
		maxRetries := 3

		for {
			if err := c.processor.Process(event); err != nil {
				log.Printf("Can't handle event: %s", err.Error())
				retries++
				if retries > maxRetries {
					log.Println("Reached maximum retries, saving event to backup.")
					if err := c.backupEvent(event); err != nil {
						log.Printf("Failed to backup event: %s", err.Error())
					}
					break
				}
				log.Printf("Retrying event processing...")
				time.Sleep(time.Duration(retries) * time.Second)
				continue
			}
			break
		}
	}

	return nil
}

func (c *Consumer) backupEvent(event events.Event) error {

	if c.backupStorage == nil {
		c.backupStorage = make(map[string]events.Event)
	}

	key := fmt.Sprintf("%s_%d", event.Text, time.Now().UnixNano())

	c.backupStorage[key] = event

	log.Printf("Event '%s' backed up with key: %s", event.Text, key)

	return nil
}

func (c *Consumer) getEventFromBackup(key string) (events.Event, error) {
	if c.backupStorage == nil {
		return events.Event{}, errors.New("No events in backup storage")
	}

	event, ok := c.backupStorage[key]
	if !ok {
		return events.Event{}, fmt.Errorf("Event with key %s not found in backup storage", key)
	}

	log.Printf("Event '%s' retrieved from backup with key: %s", event.Text, key)

	return event, nil
}
