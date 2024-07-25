package common

type (
	UserTopicRepositoryNotExistError      string
	UserTopicRepositoryAlreadyExistsError string
	UserNotExistError                     string
	UserAlreadyExistsError                string
	EmptyUserName                         string
	TopicTitleError                       string
	TopicNotExistsError                   string
	InvalidAuthData                       string
	InvalidToken                          string
)

func (e TopicTitleError) Error() string {
	return string(e)
}

func (e TopicNotExistsError) Error() string {
	return string(e)
}

func (e UserNotExistError) Error() string {
	return string(e)
}

func (e UserAlreadyExistsError) Error() string {
	return string(e)
}

func (e EmptyUserName) Error() string {
	return string(e)
}

func (e UserTopicRepositoryNotExistError) Error() string {
	return string(e)
}

func (e UserTopicRepositoryAlreadyExistsError) Error() string {
	return string(e)
}

func (e InvalidAuthData) Error() string {
	return string(e)
}

func (e InvalidToken) Error() string {
	return string(e)
}
