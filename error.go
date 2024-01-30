package main

type UserError struct {
	message string
}

func (error UserError) Error() string {
	return error.message
}
