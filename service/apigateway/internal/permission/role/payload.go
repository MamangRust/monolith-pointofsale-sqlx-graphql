package rolepermission

// RoleRequestPayload represents the base payload structure for role management requests.
// Contains common fields used across role-related operations.
type RoleRequestPayload struct {
	UserID        int    `json:"user_id"`        // ID of the user performing the role operation
	CorrelationID string `json:"correlation_id"` // Unique identifier for request tracing
	ReplyTopic    string `json:"reply_topic"`    // Topic name for asynchronous response delivery
}

// RoleResponsePayload represents the base payload structure for role validation responses.
// Used to verify role assignments and permissions.
type RoleResponsePayload struct {
	CorrelationID string   `json:"correlation_id"` // Unique ID for request tracing
	Valid         bool     `json:"valid"`          // Indicates if role validation succeeded
	RoleNames     []string `json:"role_names"`     // List of role names associated with the validation
}
