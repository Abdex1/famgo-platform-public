// STEP 9 — CREATE DLQ GOVERNANCE

// packages/kafka-sdk/consumer/dlq_handler.go


// ENTERPRISE DLQ RULE

// Every critical topic MUST have:

// <topic>.dlq
// Examples:

// ride.created.v1.dlq
// payment.completed.v1.dlq


// IMPLEMENTATION

package consumer

func BuildDLQTopic(topic string) string {
	return topic + ".dlq"
}
