/*
 * Managed Service API
 *
 * Managed Service API
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// ConnectorClusterUpdateStatusRequest Schema for the request to update a data plane cluster's status
type ConnectorClusterUpdateStatusRequest struct {
	// The cluster data plane conditions
	Conditions []DataPlaneClusterUpdateStatusRequestConditions `json:"conditions,omitempty"`
	Total      DataPlaneClusterUpdateStatusRequestTotal        `json:"total,omitempty"`
	Remaining  DataPlaneClusterUpdateStatusRequestTotal        `json:"remaining,omitempty"`
	NodeInfo   DataPlaneClusterUpdateStatusRequestNodeInfo     `json:"nodeInfo,omitempty"`
	ResizeInfo DataPlaneClusterUpdateStatusRequestResizeInfo   `json:"resizeInfo,omitempty"`
}