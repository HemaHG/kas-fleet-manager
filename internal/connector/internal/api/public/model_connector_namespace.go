/*
 * Connector Service Fleet Manager
 *
 * Connector Service Fleet Manager is a Rest API to manage connectors.
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package public

import (
	"time"
)

// ConnectorNamespace A connector namespace
type ConnectorNamespace struct {
	Id          string                                     `json:"id"`
	Kind        string                                     `json:"kind,omitempty"`
	Href        string                                     `json:"href,omitempty"`
	Owner       string                                     `json:"owner,omitempty"`
	CreatedAt   time.Time                                  `json:"created_at,omitempty"`
	ModifiedAt  time.Time                                  `json:"modified_at,omitempty"`
	Name        string                                     `json:"name"`
	Annotations []ConnectorNamespaceRequestMetaAnnotations `json:"annotations,omitempty"`
	Version     int64                                      `json:"version"`
	ClusterId   string                                     `json:"cluster_id"`
	Expiration  string                                     `json:"expiration,omitempty"`
	Tenant      ConnectorNamespaceTenant                   `json:"tenant"`
}
