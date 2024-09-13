package types

import (
	"time"

	"github.com/google/uuid"
)

// UserActivityLog represents a log entry for a user's activity
type UserActivityLog struct {
	ID        uuid.UUID              `bson:"_id" json:"id"`                                // Unique identifier for the log entry
	UserID    uuid.UUID              `bson:"user_id" json:"user_id"`                       // ID of the user who performed the activity
	Action    string                 `bson:"action" json:"action"`                         // Description of the activity performed
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`                   // When the activity occurred
	IPAddress string                 `bson:"ip_address,omitempty" json:"ip_address"`       // Optional: User's IP address
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // Optional: Additional metadata
}
