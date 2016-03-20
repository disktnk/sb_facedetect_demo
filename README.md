# face detection demo on SensorBee

With OpenCV cascade classifier function. This demo is putting a mask on faces. Target OS is Linux or Mac.

## put mask on faces demo

### prepare target file and model file

Prepare target video file, mask file and cascade model file. Model files, for example, [opencv/data](https://github.com/Itseez/opencv/tree/master/data) can use.

Below tree is a sample for data set

```
./data
├ target_video_file.avi
├ haarcascade_frontalcatface.xml
└ mask.png
```

### edit BQL

```bash
cp rectangle_mount.bql.sample rectangle_mount.bql
```

Edit file paths at rectangle_mount.bql L2, L7, L36.

```SQL
CREATE PAUSED SOURCE camera TYPE opencv_capture_from_uri WITH
    uri="data/target_video_file.avi",
    frame_skip=0, next_frame_error=false
;

CREATE STATE face_detector TYPE opencv_cascade_classifier WITH
    file="data/haarcascade_frontalcatface.xml"
;

...

CREATE STATE mount TYPE opencv_shared_image WITH
    file="data/mask.png"
;
```

### build SensorBee and make working directory

```bash
$ build_sensorbee --download-plugins=false
$ mkdir result
```

Will be made "sensorbee_main.go" and "sensorbee" (binary file). If fail to build by not fund libraries, download each library by `go get`.


### execute the BQL

```bash
$ ./sensorbee runfile -c conf.yaml rectangle_mount.bql
```

Will be created "detected.avi" and "mounted.avi" on "result" directory.


## put mask on a certain person

TODO
