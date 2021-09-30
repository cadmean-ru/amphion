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
	sceneGrid.AddRow(WrapContent)
	sceneGrid.AddRow(FillParent)
	scene.AddComponent(sceneGrid)

	padding := NewSceneObject("padding")
	padding.Transform.SetSize(0, 50)
	scene.AddChild(padding)

	buttonBar := NewSceneObject("button bar")
	buttonBarGrid := NewGridLayout()
	buttonBarGrid.Orientation = GridHorizontal
	buttonBar.AddComponent(buttonBarGrid)
	scene.AddChild(buttonBar)

	image := NewSceneObject("image")
	imageView := NewImageView("https://c.tenor.com/NTGHJ9gpDGkAAAAC/enjoying-the-ride-grogu.gif")
	image.AddComponent(imageView)
	scene.AddChild(image)

	handlePlayClick := func(event AmphionEvent) bool {
		imageView.StartAnimation()
		return true
	}

	playButton := NewSceneObject("play button")
	playButtonView := NewShapeView(ShapeRectangle)
	playButtonView.CornerRadius = 5
	playButtonView.FillColor = Green()
	playButton.AddComponent(playButtonView)
	playButton.AddComponent(NewRectBoundary())
	playButton.AddComponent(NewEventListener(EventMouseDown, handlePlayClick))
	playButton.AddComponent(NewEventListener(EventTouchDown, handlePlayClick))
	playButtonText := NewSceneObject("play button text")
	playButtonTextView := NewTextView("Play")
	playButtonText.AddComponent(playButtonTextView)
	playButton.AddChild(playButtonText)
	buttonBar.AddChild(playButton)

	handlePauseClick := func(event AmphionEvent) bool {
		imageView.StopAnimation()
		return true
	}

	pauseButton := NewSceneObject("pause button")
	pauseButtonView := NewShapeView(ShapeRectangle)
	pauseButtonView.CornerRadius = 5
	pauseButtonView.FillColor = Red()
	pauseButton.AddComponent(pauseButtonView)
	pauseButton.AddComponent(NewRectBoundary())
	pauseButton.AddComponent(NewEventListener(EventMouseDown, handlePauseClick))
	pauseButton.AddComponent(NewEventListener(EventTouchDown, handlePauseClick))
	pauseButtonText := NewSceneObject("pause button text")
	pauseButtonTextView := NewTextView("Pause")
	pauseButtonText.AddComponent(pauseButtonTextView)
	pauseButton.AddChild(pauseButtonText)
	buttonBar.AddChild(pauseButton)

	handleResumeClick := func(event AmphionEvent) bool {
		imageView.ResumeAnimation()
		return true
	}

	resumeButton := NewSceneObject("resume button")
	resumeButtonView := NewShapeView(ShapeRectangle)
	resumeButtonView.CornerRadius = 5
	resumeButtonView.FillColor = Blue()
	resumeButton.AddComponent(resumeButtonView)
	resumeButton.AddComponent(NewRectBoundary())
	resumeButton.AddComponent(NewEventListener(EventMouseDown, handleResumeClick))
	resumeButton.AddComponent(NewEventListener(EventTouchDown, handleResumeClick))
	resumeButtonText := NewSceneObject("resume button text")
	resumeButtonTextView := NewTextView("Resume")
	resumeButtonTextView.SetTextColor(White())
	resumeButtonText.AddComponent(resumeButtonTextView)
	resumeButton.AddChild(resumeButtonText)
	buttonBar.AddChild(resumeButton)

	handleChangeClick := func(event AmphionEvent) bool {
		imageView.SetResId(Res_images_giphy)
		return true
	}

	changeButton := NewSceneObject("change button")
	changeButtonView := NewShapeView(ShapeRectangle)
	changeButtonView.CornerRadius = 5
	changeButtonView.FillColor = Pink()
	changeButton.AddComponent(changeButtonView)
	changeButton.AddComponent(NewRectBoundary())
	changeButton.AddComponent(NewEventListener(EventMouseDown, handleChangeClick))
	changeButton.AddComponent(NewEventListener(EventTouchDown, handleChangeClick))
	changeButtonText := NewSceneObject("change button text")
	changeButtonTextView := NewTextView("Change")
	changeButtonText.AddComponent(changeButtonTextView)
	changeButton.AddChild(changeButtonText)
	buttonBar.AddChild(changeButton)

	handleRemoveClick := func(event AmphionEvent) bool {
		scene.RemoveChild(image)
		return true
	}

	removeButton := NewSceneObject("remove button")
	removeButtonView := NewShapeView(ShapeRectangle)
	removeButtonView.CornerRadius = 5
	removeButtonView.FillColor = Black()
	removeButton.AddComponent(removeButtonView)
	removeButton.AddComponent(NewRectBoundary())
	removeButton.AddComponent(NewEventListener(EventMouseDown, handleRemoveClick))
	removeButton.AddComponent(NewEventListener(EventTouchDown, handleRemoveClick))
	removeButtonText := NewSceneObject("remove button text")
	removeButtonTextView := NewTextView("Remove")
	removeButtonTextView.SetTextColor(White())
	removeButtonText.AddComponent(removeButtonTextView)
	removeButton.AddChild(removeButtonText)
	buttonBar.AddChild(removeButton)

	return scene
}