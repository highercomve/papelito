package helpermodels

import "time"

// Timeable interface of something with TimeStamp
type Timeable interface {
	SetCreatedAt()
	SetUpdatedAt()
	SetDeletedAt()
	GetCreatedAt() *time.Time
	GetUpdatedAt() *time.Time
	GetDeletedAt() *time.Time
}

// Timestamp timestamp extension for all models
type Timestamp struct {
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at"`
	CreatedAt *time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

// NewTimeStamp create new timestamp
func NewTimeStamp() Timestamp {
	timeNow := time.Now()
	return Timestamp{
		CreatedAt: &timeNow,
		UpdatedAt: nil,
		DeletedAt: nil,
	}
}

// SetUpdatedAt set update to a timestamp
func (t *Timestamp) SetUpdatedAt() {
	timeNow := time.Now()
	t.UpdatedAt = &timeNow
}

// SetDeletedAt set delete to a timestamp
func (t *Timestamp) SetDeletedAt() {
	if t.DeletedAt == nil {
		timeNow := time.Now()
		t.DeletedAt = &timeNow
	}
}

// SetCreatedAt to a timestamp
func (t *Timestamp) SetCreatedAt() {
	if t.CreatedAt == nil {
		timeNow := time.Now()
		t.CreatedAt = &timeNow
	}
}

// GetCreatedAt get method for created at
func (t *Timestamp) GetCreatedAt() *time.Time {
	return t.CreatedAt
}

// GetUpdatedAt get method for updated at
func (t *Timestamp) GetUpdatedAt() *time.Time {
	return t.UpdatedAt
}

// GetDeletedAt get method for deleted at
func (t *Timestamp) GetDeletedAt() *time.Time {
	return t.DeletedAt
}
