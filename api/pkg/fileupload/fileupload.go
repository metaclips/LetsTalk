package fileupload

import (
	"github.com/pion/rtp"
)

// maxDropboxUpload is 2MiB
// ToDo: all should be defined config.json
const (
	maxDropboxUpload = 2 * 1024 * 1024
	fileUploadPath   = "/class_sessions/"
)

type FileUpload interface {
	UploadFile() error
	AddVideoChunk(packet *rtp.Packet)
	AddAudioChunk(packet *rtp.Packet)
	GetUploadLink() (string, error)
	Close() error
}
