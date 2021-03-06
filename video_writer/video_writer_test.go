package videowriter

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io/ioutil"
	"os"
	"testing"
)

func TestVideoWriterCreatorCreatesSink(t *testing.T) {
	ctx := &core.Context{}
	ioParams := &bql.IOParams{}
	removeTestAVIFile()
	Convey("Given a VideoWriter sink creator", t, func() {
		vc := VideoWiterCreator{}
		params := data.Map{}
		Convey("When parameters are empty", func() {
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter not have file name", func() {
			params["fps"] = data.Float(5)
			params["width"] = data.Int(1920)
			params["height"] = data.Int(1480)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		SkipConvey("When parameter have invalid file name", func() {
			// file name will cast string use `data.ToString`, so any type can
			// cast, and this test is skipped
			params["file_name"] = data.Null{}
			params["fps"] = data.Float(5)
			params["width"] = data.Int(1920)
			params["height"] = data.Int(1480)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter have only height", func() {
			params["file_name"] = data.String("dummy")
			params["height"] = data.Int(1480)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter have only width", func() {
			params["file_name"] = data.String("dummy")
			params["width"] = data.Int(1920)
			Convey("Then sink should not be created", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldNotBeNil)
				So(sink, ShouldBeNil)
			})
		})
		Convey("When parameter have invalid values", func() {
			params["file_name"] = data.String("dummy")
			testMap := data.Map{
				"fps":    data.String("５"),
				"width":  data.String("a"),
				"height": data.String("$"),
			}
			for k, v := range testMap {
				v := v
				msg := fmt.Sprintf("%v error", k)
				Convey("Then sink should not be created because of "+msg, func() {
					params[k] = v
					s, err := vc.CreateSink(ctx, ioParams, params)
					So(err, ShouldNotBeNil)
					So(s, ShouldBeNil)
				})
			}
		})

		Convey("When parameters have only file name", func() {
			params["file_name"] = data.String("dummy")
			Convey("Then sink should be created and other parameters are set default", func() {
				sink, err := vc.CreateSink(ctx, ioParams, params)
				So(err, ShouldBeNil)
				Reset(func() {
					removeTestAVIFile()
					sink.Close(ctx)
				})
				vs, ok := sink.(*videoWriterSink)
				So(ok, ShouldBeTrue)
				So(vs.vw, ShouldNotBeNil)
				_, err = os.Stat("dummy.avi")
				So(os.IsNotExist(err), ShouldBeTrue)

				Convey("And when the sink is written an image", func() {
					img, err := ioutil.ReadFile("test_cvmat")
					So(err, ShouldBeNil)
					tu := &core.Tuple{
						Data: data.Map{
							"format": data.String("cvmat"),
							"width":  data.Int(640),
							"height": data.Int(480),
							"image":  data.Blob(img),
						},
					}
					err = sink.Write(ctx, tu)
					So(err, ShouldBeNil)
					Convey("Then should create dummy.avi", func() {
						_, err = os.Stat("dummy.avi")
						So(os.IsNotExist(err), ShouldBeFalse)

						Convey("And when another tuple is written", func() {
							err = sink.Write(ctx, tu)
							So(err, ShouldBeNil)
						})
					})
				})
			})
		})
	})
}

func removeTestAVIFile() {
	_, err := os.Stat("dummy.avi")
	if !os.IsNotExist(err) {
		os.Remove("dummy.avi")
	}
}
