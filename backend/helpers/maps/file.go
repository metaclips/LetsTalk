package maps

import "github.com/metaclips/LetsTalk/backend/helpers/constants"

var FileTypes = map[string]string{
	constants.FileTypeFile:             constants.FileTypeFile,
	constants.FileTypeMessage:          constants.FileTypeMessage,
	constants.FileTypeInfo:             constants.FileTypeInfo,
	constants.FileTypeClassSession:     constants.FileTypeClassSession,
	constants.FileTypeClassSessionLink: constants.FileTypeClassSessionLink,
}
