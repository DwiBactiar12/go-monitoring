package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"monitoring/config"
	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTClient struct {
	client         mqtt.Client
	monitoringRepo iface.MonitoringRepository
	topic          string
}

func NewMQTTClient(cfg *config.Config, repo iface.MonitoringRepository) *MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", cfg.MQTT.Broker, cfg.MQTT.Port))
	opts.SetClientID("iot_monitoring_server")

	if cfg.MQTT.Username != "" {
		opts.SetUsername(cfg.MQTT.Username)
		opts.SetPassword(cfg.MQTT.Password)
	}

	client := mqtt.NewClient(opts)

	return &MQTTClient{
		client:         client,
		topic:          cfg.MQTT.Topic,
		monitoringRepo: repo,
	}
}

func (m *MQTTClient) Start() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("Failed to connect to MQTT broker: %v", token.Error())
		return
	}

	log.Println("Connected to MQTT broker")

	// Subscribe to telemetry topic
	if token := m.client.Subscribe(m.topic+"/+/telemetry", 1, m.handleTelemetryMessage); token.Wait() && token.Error() != nil {
		log.Printf("Failed to subscribe to MQTT topic: %v", token.Error())
		return
	}

	log.Printf("Subscribed to MQTT topic: %s/+/telemetry", m.topic)
}

func (m *MQTTClient) handleTelemetryMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received MQTT message: %s", msg.Payload())

	var telemetry entity.MonitoringData
	if err := json.Unmarshal(msg.Payload(), &telemetry); err != nil {
		log.Printf("Failed to parse telemetry data: %v", err)
		return
	}

	// 	// Extract device ID from topic (topic format: iot/monitoring/{device_id}/telemetry)
	topicParts := strings.Split(msg.Topic(), "/")
	if len(topicParts) >= 3 {
		id, err := uuid.Parse(topicParts[2])
		if err != nil {
			log.Printf("Invalid UUID in topic: %v", err)
			return
		}
		telemetry.DeviceID = id
	}

	telemetry.Timestamp = time.Now()
	ctx := context.Background()
	if err := m.monitoringRepo.Store(ctx, &telemetry); err != nil {
		log.Printf("Failed to save telemetry data: %v", err)
	}
}

func (m *MQTTClient) PublishTelemetry(topic string, payload []byte) error {
	token := m.client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}
