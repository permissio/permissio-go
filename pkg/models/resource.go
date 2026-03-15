package models

import (
	"encoding/json"
	"strings"
)

// ResourceCreate represents the data for creating a new resource type.
type ResourceCreate struct {
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Actions     []string               `json:"-"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// MarshalJSON serializes ResourceCreate, converting Actions []string into the
// map format the backend expects: {"read": {"name": "Read"}, ...}
func (r ResourceCreate) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Key         string                 `json:"key"`
		Name        string                 `json:"name,omitempty"`
		Description string                 `json:"description,omitempty"`
		Actions     map[string]interface{} `json:"actions,omitempty"`
		Attributes  map[string]interface{} `json:"attributes,omitempty"`
	}

	a := Alias{
		Key:         r.Key,
		Name:        r.Name,
		Description: r.Description,
		Attributes:  r.Attributes,
	}

	if len(r.Actions) > 0 {
		actionsMap := make(map[string]interface{}, len(r.Actions))
		for _, action := range r.Actions {
			actionsMap[action] = map[string]interface{}{"name": capitalizeFirst(action)}
		}
		a.Actions = actionsMap
	}

	return json.Marshal(a)
}

// capitalizeFirst returns the string with the first letter uppercased.
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// NewResourceCreate creates a new ResourceCreate with the given key.
func NewResourceCreate(key string) *ResourceCreate {
	return &ResourceCreate{Key: key}
}

// SetName sets the name for the resource.
func (r *ResourceCreate) SetName(name string) *ResourceCreate {
	r.Name = name
	return r
}

// SetDescription sets the description for the resource.
func (r *ResourceCreate) SetDescription(description string) *ResourceCreate {
	r.Description = description
	return r
}

// SetActions sets the actions for the resource.
func (r *ResourceCreate) SetActions(actions []string) *ResourceCreate {
	r.Actions = actions
	return r
}

// AddAction adds an action to the resource.
func (r *ResourceCreate) AddAction(action string) *ResourceCreate {
	r.Actions = append(r.Actions, action)
	return r
}

// SetAttributes sets custom attributes for the resource.
func (r *ResourceCreate) SetAttributes(attributes map[string]interface{}) *ResourceCreate {
	r.Attributes = attributes
	return r
}

// ResourceUpdate represents the data for updating a resource.
type ResourceUpdate struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Actions     []string               `json:"-"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// MarshalJSON serializes ResourceUpdate, converting Actions []string into the
// map format the backend expects: {"read": {"name": "Read"}, ...}
func (r ResourceUpdate) MarshalJSON() ([]byte, error) {
	type Alias struct {
		Name        *string                `json:"name,omitempty"`
		Description *string                `json:"description,omitempty"`
		Actions     map[string]interface{} `json:"actions,omitempty"`
		Attributes  map[string]interface{} `json:"attributes,omitempty"`
	}

	a := Alias{
		Name:        r.Name,
		Description: r.Description,
		Attributes:  r.Attributes,
	}

	if len(r.Actions) > 0 {
		actionsMap := make(map[string]interface{}, len(r.Actions))
		for _, action := range r.Actions {
			actionsMap[action] = map[string]interface{}{"name": capitalizeFirst(action)}
		}
		a.Actions = actionsMap
	}

	return json.Marshal(a)
}

// ResourceRead represents a resource returned from the API.
type ResourceRead struct {
	ID          string                 `json:"id"`
	Key         string                 `json:"key"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Actions     []string               `json:"-"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

// UnmarshalJSON deserializes ResourceRead, converting the backend's actions map
// format {"read": {...}, "write": {...}} back to a []string slice.
func (r *ResourceRead) UnmarshalJSON(data []byte) error {
	type Alias struct {
		ID          string                 `json:"id"`
		Key         string                 `json:"key"`
		Name        string                 `json:"name,omitempty"`
		Description string                 `json:"description,omitempty"`
		Actions     interface{}            `json:"actions,omitempty"`
		Attributes  map[string]interface{} `json:"attributes,omitempty"`
		CreatedAt   string                 `json:"created_at"`
		UpdatedAt   string                 `json:"updated_at"`
	}

	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	r.ID = a.ID
	r.Key = a.Key
	r.Name = a.Name
	r.Description = a.Description
	r.Attributes = a.Attributes
	r.CreatedAt = a.CreatedAt
	r.UpdatedAt = a.UpdatedAt

	// Actions can be a map (from backend) or a slice (older format)
	switch v := a.Actions.(type) {
	case map[string]interface{}:
		actions := make([]string, 0, len(v))
		for k := range v {
			actions = append(actions, k)
		}
		r.Actions = actions
	case []interface{}:
		actions := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				actions = append(actions, s)
			}
		}
		r.Actions = actions
	default:
		r.Actions = nil
	}

	return nil
}

// ResourceList represents a paginated list of resources.
type ResourceList struct {
	Data []ResourceRead `json:"data"`
	PaginatedResponse
}

// ResourceListParams represents parameters for listing resources.
type ResourceListParams struct {
	ListParams
	Search string `json:"search,omitempty"`
}

// ResourceInstanceCreate represents the data for creating a resource instance.
type ResourceInstanceCreate struct {
	Key          string                 `json:"key"`
	ResourceType string                 `json:"resource_type"`
	Tenant       string                 `json:"tenant,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// ResourceInstanceRead represents a resource instance returned from the API.
type ResourceInstanceRead struct {
	ID           string                 `json:"id"`
	Key          string                 `json:"key"`
	ResourceType string                 `json:"resource_type"`
	Tenant       string                 `json:"tenant,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}
