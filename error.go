package main

type ReasonType int16

const (
	UsernameAlreadyTaken ReasonType = iota
	EmailAlreadyRegistered
	UserNotRegistered
	LoginAttemptFailed
)

type DomainError interface {
	getReason() ReasonType
}

type UserError struct {
	reason  ReasonType
	message string
}

type SystemError struct {
	error
}

func (error UserError) Error() string {
	return error.message
}

func (error UserError) getReason() ReasonType {
	return error.reason
}

func IsAuthenticationErrorReason(err error) bool {
	if userErr, ok := err.(UserError); ok {
		if userErr.reason == LoginAttemptFailed {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func IsUserNotRegistedError(err)
