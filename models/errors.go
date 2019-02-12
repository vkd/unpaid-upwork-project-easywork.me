package models

type TrackerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e TrackerError) Error() string {
	return e.Message
}

var IncorrectJSONInputError = TrackerError{1000, "Wrong JSON input"}
var WrongEmailOrPassword = TrackerError{1010, "Wrong email or password"}
var AccessForbidden = TrackerError{1015, "Access to the resource is forbidden for current user"}

var JsonDecodeError = TrackerError{1020, "Request Json Decoding Error"}
var PayloadValidationError = TrackerError{1030, "Payload Validation Error"}
var JwtTokenParseError = TrackerError{1040, "JWT token parsing Error"}
var JwtTokenExpiredError = TrackerError{1050, "JWT token Expired"}
var JwtClaimsError = TrackerError{1060, "JWT claims Error"}

var UserEmailExists = TrackerError{2080, "User email already exists and can't be used for a new registration"}
var WrongUsername = TrackerError{2081, "Username can only contain alphanumeric characters and be in lowercase"}
var EmptyFirstName = TrackerError{2081, "First name must exist"}
var EmptyLastName = TrackerError{2081, "Last name must exist"}
var UserEmailOrPasswordEmpty = TrackerError{2090, "User email or password for registration is empty"}

var UserNotFound = TrackerError{2091, "User not found"}
var InviteeUserNotFound = TrackerError{2091, "Invitee User not found"}
var ProjectNotFound = TrackerError{2092, "Project not found"}
var InvitationNotFound = TrackerError{2093, "Invitation not found"}

var ContractNotFound = TrackerError{2094, "Contract not found"}
var ContractAlreadyStarted = TrackerError{2094, "Contract already started"}
var ContractAlreadyPaused = TrackerError{2094, "Contract already paused"}
var ContractAlreadySameStatus = TrackerError{2094, "Contract already has same status"}
var ContractNotChangedStatus = TrackerError{2094, "Contract's status not changed"}

var TermsNotFound = TrackerError{2095, "Terms not found"}

var ProjectDoesntBelongToUser = TrackerError{2096, "User doesn't have a project with provided id"}

var UserCannotBeInvitedToHisOwnProject = TrackerError{2097, "User cannot invite himself to his own project"}

var ContractLogAfterStopError = TrackerError{2097, `"log" event cannot be saved when last event was "stop"`}
var ContractStartAfterStartError = TrackerError{2097, `"start" event cannot be saved when last event was "start"`}
var ContractStopAfterStopError = TrackerError{2097, `"stop" event cannot be saved when last event was "stop"`}
var ContractStartAfterLogError = TrackerError{2097, `"start" event cannot be saved when last event was "log"`}

var ContractIsNotStarted = TrackerError{2097, `Contract is not in "started" state and events cannot be logged into it`}
