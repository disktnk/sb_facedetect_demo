package videowriter

import (
	"fmt"
	"gopkg.in/sensorbee/opencv.v0"
	"gopkg.in/sensorbee/opencv.v0/bridge"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"os"
	"path/filepath"
	"sync"
)

var (
	fileNamePath = data.MustCompilePath("file_name")
	fpsPath      = data.MustCompilePath("fps")
	widthPath    = data.MustCompilePath("width")
	heightPath   = data.MustCompilePath("height")
)

// VideoWiterCreator is a creator of VideoWriter.
type VideoWiterCreator struct{}

// CreateSink creates a AVI Video writer sink, which outputs video file with
// input image data.
//
// WITH parameters.
//
// file_name: [required] AVI filename, will be created [file_name].avi.
//
// fps: FPS, if empty then set 1.0.
//
// width: Width of output video file, if empty then the video writer will
// initialize with the first image.
//
// height: Height of output video file, if empty then the video writer will
// initialize with the first image
func (c *VideoWiterCreator) CreateSink(ctx *core.Context, ioParams *bql.IOParams,
	params data.Map) (core.Sink, error) {

	fn, err := params.Get(fileNamePath)
	if err != nil {
		return nil, err
	}
	name, err := data.ToString(fn)
	if err != nil {
		return nil, err
	}
	name += ".avi"

	absPath, err := filepath.Abs(name)
	if err != nil {
		return nil, fmt.Errorf("invalid file path: %v", err.Error())
	}
	dirPath := filepath.Dir(absPath)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		os.MkdirAll(dirPath, 0755)
	}

	fps := float64(1.0)
	if f, err := params.Get(fpsPath); err == nil {
		if fps, err = data.ToFloat(f); err != nil {
			return nil, err
		}
	}

	width := int64(0)
	widthFlag := false
	height := int64(0)
	heightFlag := false
	if w, err := params.Get(widthPath); err == nil {
		if width, err = data.ToInt(w); err != nil {
			return nil, err
		}
		widthFlag = true
	}
	if h, err := params.Get(heightPath); err == nil {
		if height, err = data.ToInt(h); err != nil {
			return nil, err
		}
		heightFlag = true
	}
	if widthFlag != heightFlag {
		return nil, fmt.Errorf("both width and height must be set up")
	}

	return &videoWriterSink{
		name:   name,
		fps:    fps,
		width:  int(width),
		height: int(height),
		vw:     bridge.NewVideoWriter(),
	}, nil
}

type videoWriterSink struct {
	name   string
	fps    float64
	width  int
	height int
	vw     bridge.VideoWriter
	mu     sync.RWMutex
}

// Write input images and add to a video file which have been created when the
// sink is created. Input tuples are expected to have `opencv.RawData`
// structure which include "width", "height" and "image".
func (s *videoWriterSink) Write(ctx *core.Context, t *core.Tuple) error {
	raw, err := opencv.ConvertMapToRawData(t.Data)
	if err != nil {
		return err
	}
	imgp, err := raw.ToMatVec3b()
	if err != nil {
		return err
	}
	defer imgp.Delete()

	if !s.vw.IsOpened() {
		if err := s.open(imgp); err != nil {
			return err
		}
	}

	s.vw.Write(imgp)
	return nil
}

func (s *videoWriterSink) open(img bridge.MatVec3b) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.width != 0 {
		s.vw.Open(s.name, s.fps, s.width, s.height)
	} else {
		s.vw.OpenWithMat(s.name, s.fps, img)
	}
	if !s.vw.IsOpened() {
		return fmt.Errorf("cannot video writer open: %v", s.name)
	}
	return nil
}

func (s *videoWriterSink) Close(ctx *core.Context) error {
	s.vw.Delete()
	return nil
}
