package rabbitmq

type MQMessage struct {
	EmailTo string `json:"emailto"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
