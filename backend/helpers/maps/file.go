package maps

import "github.com/metaclips/LetsTalk/backend/helpers/constants"

var MessageTypes = map[string]string{
	constants.MessageTypeFile:             constants.MessageTypeFile,
	constants.MessageTypeMessage:          constants.MessageTypeMessage,
	constants.MessageTypeInfo:             constants.MessageTypeInfo,
	constants.MessageTypeClassSession:     constants.MessageTypeClassSession,
	constants.MessageTypeClassSessionLink: constants.MessageTypeClassSessionLink,
}
