/*
 * Connector Service Fleet Manager Private APIs
 *
 * Connector Service Fleet Manager apis that are used by internal services.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

// KafkaConnectionSettings Holds the configuration to connect to a Kafka Instance.
type KafkaConnectionSettings struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}
