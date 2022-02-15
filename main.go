//go:build (windows || darwin || linux) && !android && !ios
// +build windows darwin linux
// +build !android
// +build !ios

package main

import (
	"flag"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/cadmean-ru/amphion/scenes"
	"github.com/cadmean-ru/amphion/utils"
	"log"
	"runtime"
)

//go:generate ./build/darwin/test --generate resources -i ./build/darwin/res -o ./res.gen.go

func init() {
	runtime.LockOSThread()
}

func main() {

	//region generator
	var generator string
	var inputPath string
	var outputPath string
	var packageName string

	flag.StringVar(&generator, "generate", "", "Launches the specified code generator instead of starting engine. "+
		"Available generators are: shaders.")
	flag.StringVar(&inputPath, "i", "", "Define input path")
	flag.StringVar(&outputPath, "o", "", "Define output path")
	flag.StringVar(&packageName, "package", "", "Define package name")
	flag.Parse()

	switch generator {
	case "shaders":
		utils.GenerateShaders(inputPath, outputPath, packageName)
		return
	case "resources":
		utils.GenerateResourcesIndices(inputPath, outputPath)
		return
	default:
		break
	}

	//endregion

	front := pc.NewFrontend()
	front.Init()

	e := engine.Initialize(front)

	scenes.RegisterComponents(e)
	scenes.RegisterResources(e)

	go func() {
		e.Start()

		if err := e.ShowScene(scenes.MarginScene(e)); err != nil {
			log.Println(err)
		}

		// e.LoadApp()
	}()

	front.Run()
}
