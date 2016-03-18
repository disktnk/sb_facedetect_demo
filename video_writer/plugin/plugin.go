package plugin

import (
	"github.com/disktnk/sb_facedetect_demo/video_writer"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
)

func init() {
	bql.MustRegisterGlobalSinkCreator("sbfddemo_avi_writer",
		&videowriter.VideoWiterCreator{})
}
