package message

type EventMessage interface {
	Message() ([]byte, error)
}

type CommentMessage interface {
	SendUpdateCommentCountEvent(message EventMessage) error
}
