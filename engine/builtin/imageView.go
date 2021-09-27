package builtin

import (
	"bytes"
	"errors"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"image"
	"image/draw"
	"image/gif"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
)

//ImageView displays an image given its resource id or url.
//Supports jpeg, png and gif formats.
//Supports file system path or http urls.
type ImageView struct {
	engine.ViewImpl
	ResId         a.ResId `state:"resId"`
	ImageUrl      string  `state:"imageUrl"`
	imageType     string
	frames        []*rendering.Bitmap
	framesCount   int
	currentFrame  int
	delays        []int
	nextFrameTime time.Time
	loopCount     int
	currentLoop   int
}

func (v *ImageView) OnStart() {
	v.loadImageTask()
}

func (v *ImageView) loadImageTask() {
	task := engine.NewTaskBuilder().
		Run(v.loadImage).
		Then(func(res interface{}) {
			if v.PrimitiveId == -1 {
				v.PrimitiveId = v.SceneObject.GetRenderingNode().AddPrimitive()
			}

			engine.RequestRendering()
		}).
		Err(func(err error) {
			panic(err)
		}).
		Build()
	engine.RunTask(task)
}

func (v *ImageView) loadImage() (interface{}, error) {
	if v.ImageUrl != "" && strings.HasPrefix(v.ImageUrl, "https://") {
		return v.downloadImage()
	}

	return v.loadLocalImage()
}

func (v *ImageView) loadLocalImage() (interface{}, error) {
	var resId a.ResId
	if v.ResId < 0 {
		resId = engine.GetResourceManager().IdOf(v.ImageUrl)
	} else {
		resId = v.ResId
	}

	imageBytes, err := engine.GetResourceManager().ReadFile(resId)
	if err != nil {
		return nil, err
	}

	err = v.decodeImage(imageBytes)
	return nil, err
}

func (v *ImageView) downloadImage() (interface{}, error) {
	response, err := http.Get(v.ImageUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = v.decodeImage(imageBytes)
	return nil, err
}

func (v *ImageView) decodeImage(imageBytes []byte) error {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/jpeg", "image/png":
		return v.decodeJpgOrPng(imageBytes)
	case "image/gif":
		return v.decodeGif(imageBytes)
	}
	return errors.New("unsupported format")
}

func (v *ImageView) decodeJpgOrPng(imageBytes []byte) error {
	img, format, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	frame := &rendering.Bitmap{
		Pixels: rgba.Pix,
		Width:  rgba.Bounds().Size().X,
		Height: rgba.Bounds().Size().Y,
	}
	v.frames = []*rendering.Bitmap{frame}
	v.imageType = format
	v.framesCount = 1
	v.currentFrame = 0
	v.loopCount = -1
	v.currentLoop = 0

	return nil
}

func (v *ImageView) decodeGif(imageBytes []byte) error {
	gifImage, err := gif.DecodeAll(bytes.NewReader(imageBytes))
	if err != nil {
		return err
	}

	v.frames = make([]*rendering.Bitmap, len(gifImage.Image))
	width, height := v.getGifDimensions(gifImage)
	originalBounds := image.Rect(0, 0, width, height)
	rgba := image.NewRGBA(originalBounds)
	for i, g := range gifImage.Image {
		draw.Draw(rgba, g.Bounds(), g, g.Bounds().Min, draw.Over)

		pix := make([]byte, len(rgba.Pix))
		copy(pix, rgba.Pix)

		v.frames[i] = &rendering.Bitmap{
			Pixels: pix,
			Width:  width,
			Height: height,
		}

		disposal := gifImage.Disposal[i]
		if disposal == 0 || disposal == gif.DisposalBackground {
			rgba = image.NewRGBA(originalBounds)
		}
	}
	v.imageType = "gif"
	v.framesCount = len(v.frames)
	v.currentFrame = 0
	v.delays = gifImage.Delay
	v.calculateNextFrameTime()
	v.loopCount = gifImage.LoopCount
	v.currentLoop = 0

	return nil
}

func (v *ImageView) getGifDimensions(gif *gif.GIF) (x, y int) {
	var lowestX int
	var lowestY int
	var highestX int
	var highestY int

	for _, img := range gif.Image {
		if img.Rect.Min.X < lowestX {
			lowestX = img.Rect.Min.X
		}
		if img.Rect.Min.Y < lowestY {
			lowestY = img.Rect.Min.Y
		}
		if img.Rect.Max.X > highestX {
			highestX = img.Rect.Max.X
		}
		if img.Rect.Max.Y > highestY {
			highestY = img.Rect.Max.Y
		}
	}

	return highestX - lowestX, highestY - lowestY
}

func (v *ImageView) calculateNextFrameTime() {
	v.nextFrameTime = time.Now().Add(time.Duration(v.delays[0]) * 10 * time.Millisecond)
}

func (v *ImageView) removeImage() {
	v.framesCount = 0
	v.frames = nil
	v.loopCount = -1
	v.currentLoop = 0
	v.currentFrame = 0
	v.SceneObject.GetRenderingNode().RemovePrimitive(v.PrimitiveId)
	v.PrimitiveId = -1
}

func (v *ImageView) OnUpdate(_ engine.UpdateContext) {
	if v.framesCount <= 1 {
		return
	}

	if v.nextFrameTime.Sub(time.Now()) <= 0 {
		v.currentFrame += 1
		if v.currentFrame >= v.framesCount {
			v.currentFrame = 0
			v.currentLoop = (v.currentLoop + 1) % math.MaxInt64
		}
		v.calculateNextFrameTime()
		v.Redraw()
		engine.RequestRendering()
	} else {
		engine.RequestUpdate()
	}
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	if v.framesCount == 0 {
		return
	}

	pr := rendering.NewImagePrimitive(v.frames, v.currentFrame)
	pr.Transform = v.SceneObject.Transform.ToRenderingTransform()
	ctx.GetRenderingNode().SetPrimitive(v.PrimitiveId, pr)
	v.ShouldRedraw = false
}

func (v *ImageView) OnStop() {
	v.removeImage()
}

func (v *ImageView) MeasureContents() a.Vector3 {
	if v.framesCount == 0 {
		return a.ZeroVector()
	}

	frame := v.frames[v.currentFrame]
	return a.NewIntVector3(frame.Width, frame.Height, 0).ToFloat()
}

// SetResId sets the resource index equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetResId(i a.ResId) {
	v.removeImage()
	v.ResId = i
	v.ImageUrl = ""
	v.ShouldRedraw = true
	v.loadImageTask()
}

// SetImageUrl sets the image url equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetImageUrl(url string) {
	v.removeImage()
	v.ResId = -1
	v.ImageUrl = url
	v.ShouldRedraw = true
	v.loadImageTask()
}

//StartAnimation starts the animation of a gif image from the beginning.
func (v *ImageView) StartAnimation() {
	if v.framesCount <= 1 || len(v.frames) == 1 {
		return
	}

	v.currentFrame = 0
	v.currentLoop = 0
	v.framesCount = len(v.frames)
	v.calculateNextFrameTime()
	engine.RequestUpdate()
}

//StopAnimation stops the animation of a gif image.
func (v *ImageView) StopAnimation() {
	if v.framesCount <= 1 {
		return
	}

	v.framesCount = 1
}

//ResumeAnimation resumes the animation of a gif image from the frame it was stopped on.
func (v *ImageView) ResumeAnimation() {
	if v.framesCount <= 1 || len(v.frames) == 1 {
		return
	}

	v.framesCount = len(v.frames)
	v.calculateNextFrameTime()
	engine.RequestUpdate()
}

//IsAnimating checks if this image is currently animating.
func (v *ImageView) IsAnimating() bool {
	return v.framesCount > 1
}

//ImageType returns the type of loaded image.
//The possible values are - "jpeg", "png" or "gif".
func (v *ImageView) ImageType() string {
	return v.imageType
}

//NewImageView creates a new ImageView given a resource id or an url of an image.
func NewImageView(image interface{}) *ImageView {
	switch image.(type) {
	case a.ResId:
		return &ImageView{
			ResId: image.(a.ResId),
		}
	case string:
		return &ImageView{
			ImageUrl: image.(string),
		}
	}

	panic("unknown argument type")
}
