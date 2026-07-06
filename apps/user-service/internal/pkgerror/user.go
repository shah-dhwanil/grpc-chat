package pkgerror

import "github.com/shah-dhwanil/grpc-chat/packages/pkgerror"

func NewUserNotFoundError(err error) error {
	return pkgerror.NewAppError(
		pkgerror.ResourceNotFound,
		"USER_NOT_FOUND",
		"User with requested data not found",
		nil,
		err,
	)
}

func NewUserAlreadyExistsError(err error) error {
	return pkgerror.NewAppError(
		pkgerror.ResourceAlreadyExists,
		"USER_ALREADY_EXISTS",
		"User with requested data already exists",
		nil,
		err,
	)
}
