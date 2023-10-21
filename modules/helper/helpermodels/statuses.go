package helpermodels

// StatusType type of status
type StatusType string

// Status status for a collection
type Status struct {
	State   StatusType `json:"state" bson:"state"`
	Message string     `json:"message" bson:"message"`
}

// Statusable make something Status able
type Statusable interface {
	SetDone()
	SetNew()
	SetScheduled()
	SetFailed()
	SetCanceled()
	SetDeleted()
	SetPaused()
	SetDeleting()
	SetActive()
	GetStatus() StatusType
	GetMessage() string
	IsFinal() bool
	IsValid() bool
}

const (

	// StatusDone the thing is done
	StatusDone = StatusType("done")

	// StatusNew the thing is just created
	StatusNew = StatusType("new")

	// StatusScheduled the thing is scheduled (has a date when it will start)
	StatusScheduled = StatusType("scheduled")

	// StatusActive means the thing is active (has been consume)
	StatusActive = StatusType("active")

	// StatusFailed the thing failed and has an error
	StatusFailed = StatusType("failed")

	// StatusCanceled the thing has been canceled for some reason
	StatusCanceled = StatusType("canceled")

	// StatusPaused the thing has been paused for some reason
	StatusPaused = StatusType("paused")

	// StatusSkipped the thing has been skipped for some reason
	StatusSkipped = StatusType("skipped")

	// StatusDeleted the thing has been deleted for some reason
	StatusDeleted = StatusType("deleted")

	// StatusDeleting something that as been marked to be deleted
	StatusDeleting = StatusType("deleting")
)

// ValidStatusTypes valid types of permissions
var ValidStatusTypes = map[StatusType]StatusType{
	StatusDone:      StatusDone,
	StatusNew:       StatusNew,
	StatusScheduled: StatusScheduled,
	StatusFailed:    StatusFailed,
	StatusCanceled:  StatusCanceled,
	StatusDeleted:   StatusDeleted,
	StatusDeleting:  StatusDeleting,
	StatusActive:    StatusActive,
}

// NewStatus create new status
func NewStatus(status string) *Status {
	return &Status{
		State: StatusType(status),
	}
}

// SetActive set status as active
func (s *Status) SetActive() {
	s.State = StatusActive
}

// SetDone set the status as done
func (s *Status) SetDone() {
	s.State = StatusDone
}

// SetNew set the status as New
func (s *Status) SetNew() {
	s.State = StatusNew
}

// SetScheduled set the status as Scheduled
func (s *Status) SetScheduled() {
	s.State = StatusScheduled
}

// SetFailed set the status as Failed
func (s *Status) SetFailed() {
	s.State = StatusFailed
}

// SetCanceled set the status as Canceled
func (s *Status) SetCanceled() {
	s.State = StatusCanceled
}

// SetDeleted set the status as Canceled
func (s *Status) SetDeleted() {
	s.State = StatusDeleted
}

// SetPaused set the status as paused
func (s *Status) SetPaused() {
	s.State = StatusPaused
}

// SetDeleting set element as deleting
func (s *Status) SetDeleting() {
	s.State = StatusDeleting
}

// GetStatus Get the status
func (s *Status) GetStatus() StatusType {
	return s.State
}

// GetMessage Get the message
func (s *Status) GetMessage() string {
	return s.Message
}

// IsFinal Status is final
func (s *Status) IsFinal() bool {
	return s.State == StatusDone ||
		s.State == StatusDeleted ||
		s.State == StatusCanceled
}

// IsValid permission is valid
func (s *Status) IsValid() bool {
	_, ok := ValidStatusTypes[s.State]
	return ok
}
