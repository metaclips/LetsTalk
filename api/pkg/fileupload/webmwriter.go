package fileupload

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/at-wat/ebml-go/webm"
	"github.com/metaclips/LetsTalk/api/pkg/dropbox"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v2/pkg/media/samplebuilder"
)

type webmWriter struct {
	dropbox                        *dropbox.Dropbox
	fileName                       string
	audioWriter, videoWriter       webm.BlockWriteCloser
	audioBuilder, videoBuilder     *samplebuilder.SampleBuilder
	audioTimestamp, videoTimestamp uint32
}

// NewWebmWriter writes video class session to either be upload to a Dropbox drive
// or if a token is not specified, saved to mongodb using gridFS.
func NewWebmWriter(fileName string) FileUpload {
	return &webmWriter{
		fileName:     fileName,
		audioBuilder: samplebuilder.New(10, &codecs.OpusPacket{}),
		videoBuilder: samplebuilder.New(10, &codecs.VP8Packet{}),
	}
}

func (s *webmWriter) initWriter(width, height int) {
	w, err := os.OpenFile(s.fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorln("error opening file", err)
	}

	ws, err := webm.NewSimpleBlockWriter(w,
		[]webm.TrackEntry{
			{
				Name:            "Audio",
				TrackNumber:     1,
				TrackUID:        12345,
				CodecID:         "A_OPUS",
				TrackType:       2,
				DefaultDuration: 20000000,
				Audio: &webm.Audio{
					SamplingFrequency: 48000.0,
					Channels:          2,
				},
			}, {
				Name:            "Video",
				TrackNumber:     2,
				TrackUID:        67890,
				CodecID:         "V_VP8",
				TrackType:       1,
				DefaultDuration: 33333333,
				Video: &webm.Video{
					PixelWidth:  uint64(width),
					PixelHeight: uint64(height),
				},
			},
		})
	if err != nil {
		log.Println("error initiating file writer", err)
	}

	log.Infof("WebM saver has started with video width=%d, height=%d\n", width, height)
	s.audioWriter = ws[0]
	s.videoWriter = ws[1]
}

// Append Audi file.
func (s *webmWriter) AddAudioChunk(rtpPacket *rtp.Packet) {
	s.audioBuilder.Push(rtpPacket)

	for {
		sample := s.audioBuilder.Pop()
		if sample == nil {
			return
		}

		if s.audioWriter != nil {
			s.audioTimestamp += sample.Samples
			t := s.audioTimestamp / 48
			if _, err := s.audioWriter.Write(true, int64(t), sample.Data); err != nil {
				log.Errorln("error writing audio byte", err)
			}
		}
	}
}

// Push to video.
func (s *webmWriter) AddVideoChunk(rtpPacket *rtp.Packet) {
	s.videoBuilder.Push(rtpPacket)

	for {
		sample := s.videoBuilder.Pop()
		if sample == nil {
			return
		}

		// Read VP8 header.
		videoKeyframe := (sample.Data[0]&0x1 == 0)
		if videoKeyframe {
			// Keyframe has frame information.
			raw := uint(sample.Data[6]) | uint(sample.Data[7])<<8 | uint(sample.Data[8])<<16 | uint(sample.Data[9])<<24
			width := int(raw & 0x3FFF)
			height := int((raw >> 16) & 0x3FFF)

			if s.videoWriter == nil || s.audioWriter == nil {
				// Initialize WebM saver using received frame size.
				s.initWriter(width, height)
			}
		}

		if s.videoWriter != nil {
			s.videoTimestamp += sample.Samples
			t := s.videoTimestamp / 90
			if _, err := s.videoWriter.Write(videoKeyframe, int64(t), sample.Data); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *webmWriter) Close() error {
	log.Info("Finalizing webm...\n")
	if s.audioWriter != nil {
		if err := s.audioWriter.Close(); err != nil {
			return err
		}
	}

	if s.videoWriter != nil {
		if err := s.videoWriter.Close(); err != nil {
			return err
		}
	}

	log.Infoln("video writer closed for session, uploading.", s.fileName)
	return nil
}

func (s *webmWriter) UploadFile() error {
	d, err := dropbox.NewDropbox(s.fileName)
	if err != nil {
		return err
	}

	s.dropbox = d

	return s.dropbox.UploadFile()
}

func (s *webmWriter) GetUploadLink() (string, error) {
	if s.dropbox == nil {
		return "", errors.New("dropbox instance not created")
	}
	return s.dropbox.GetUploadLink()
}
