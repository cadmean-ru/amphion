package scenes

import (
	. "github.com/cadmean-ru/amphion/common/a"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func GifScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("gif scene")
	sceneGrid := NewGridLayout()
	sceneGrid.AddColumn(FillParent)
	sceneGrid.AddRow(WrapContent)
	sceneGrid.AddRow(FillParent)
	scene.AddComponent(sceneGrid)

	buttonBar := NewSceneObject("button bar")
	buttonBarGrid := NewGridLayout()
	buttonBarGrid.Orientation = GridHorizontal
	buttonBar.AddComponent(buttonBarGrid)
	scene.AddChild(buttonBar)

	image := NewSceneObject("image")
	imageView := NewImageView("https://c.tenor.com/NTGHJ9gpDGkAAAAC/enjoying-the-ride-grogu.gif")
	image.AddComponent(imageView)
	scene.AddChild(image)

	playButton := NewSceneObject("play button")
	playButtonView := NewShapeView(ShapeRectangle)
	playButtonView.CornerRadius = 5
	playButtonView.FillColor = Green()
	playButton.AddComponent(playButtonView)
	playButton.AddComponent(NewRectBoundary())
	playButton.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		imageView.StartAnimation()
		return true
	}))
	playButtonText := NewSceneObject("play button text")
	playButtonTextView := NewTextView("Play")
	playButtonText.AddComponent(playButtonTextView)
	playButton.AddChild(playButtonText)
	buttonBar.AddChild(playButton)

	pauseButton := NewSceneObject("pause button")
	pauseButtonView := NewShapeView(ShapeRectangle)
	pauseButtonView.CornerRadius = 5
	pauseButtonView.FillColor = Red()
	pauseButton.AddComponent(pauseButtonView)
	pauseButton.AddComponent(NewRectBoundary())
	pauseButton.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		imageView.StopAnimation()
		return true
	}))
	pauseButtonText := NewSceneObject("pause button text")
	pauseButtonTextView := NewTextView("Pause")
	pauseButtonText.AddComponent(pauseButtonTextView)
	pauseButton.AddChild(pauseButtonText)
	buttonBar.AddChild(pauseButton)

	resumeButton := NewSceneObject("resume button")
	resumeButtonView := NewShapeView(ShapeRectangle)
	resumeButtonView.CornerRadius = 5
	resumeButtonView.FillColor = Blue()
	resumeButton.AddComponent(resumeButtonView)
	resumeButton.AddComponent(NewRectBoundary())
	resumeButton.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		imageView.ResumeAnimation()
		return true
	}))
	resumeButtonText := NewSceneObject("resume button text")
	resumeButtonTextView := NewTextView("Resume")
	resumeButtonTextView.SetTextColor(White())
	resumeButtonText.AddComponent(resumeButtonTextView)
	resumeButton.AddChild(resumeButtonText)
	buttonBar.AddChild(resumeButton)

	changeButton := NewSceneObject("change button")
	changeButtonView := NewShapeView(ShapeRectangle)
	changeButtonView.CornerRadius = 5
	changeButtonView.FillColor = Pink()
	changeButton.AddComponent(changeButtonView)
	changeButton.AddComponent(NewRectBoundary())
	changeButton.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		imageView.SetResId(Res_images_giphy)
		return true
	}))
	changeButtonText := NewSceneObject("change button text")
	changeButtonTextView := NewTextView("Change")
	changeButtonText.AddComponent(changeButtonTextView)
	changeButton.AddChild(changeButtonText)
	buttonBar.AddChild(changeButton)

	removeButton := NewSceneObject("remove button")
	removeButtonView := NewShapeView(ShapeRectangle)
	removeButtonView.CornerRadius = 5
	removeButtonView.FillColor = Black()
	removeButton.AddComponent(removeButtonView)
	removeButton.AddComponent(NewRectBoundary())
	removeButton.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		scene.RemoveChild(image)
		return true
	}))
	removeButtonText := NewSceneObject("remove button text")
	removeButtonTextView := NewTextView("Remove")
	removeButtonTextView.SetTextColor(White())
	removeButtonText.AddComponent(removeButtonTextView)
	removeButton.AddChild(removeButtonText)
	buttonBar.AddChild(removeButton)

	return scene
}
