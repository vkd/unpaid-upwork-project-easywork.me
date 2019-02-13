package httpservice

type TrackerMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *TrackerMessage) Description() string {
	return e.Message
}

var UserDeleted = TrackerMessage{1001, "User deleted successfully"}
var ProjectDeleted = TrackerMessage{1002, "Project deleted successfully"}
var InvitationDeleted = TrackerMessage{1003, "Invitation deleted successfully"}

var ContractStarted = TrackerMessage{1004, "Contract started successfully"}
var ContractPaused = TrackerMessage{1005, "Contract paused successfully"}
var ContractResumed = TrackerMessage{1004, "Contract resumed successfully"}
var ContractEnded = TrackerMessage{1006, "Contract ended successfully"}

var InvitationAccepted = TrackerMessage{1005, "Invitation accepted successfully"}
var InvitationDeclined = TrackerMessage{1006, "Invitation declined successfully"}
