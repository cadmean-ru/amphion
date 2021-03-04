//+build windows darwin linux

package main

import (
	"flag"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/pc"
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

	flag.StringVar(&generator, "generate", "", "Launches the specified code generator instead of starting engine. " +
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

	registerComponents(e)
	registerResources(e)

	//if data, err := scene.EncodeToYaml(); err == nil {
	//	fmt.Println(string(data))
	//}

	go func() {
		e.Start()

		if err := e.ShowScene(prefabScene(e)); err != nil {
			log.Println(err)
		}

		//e.LoadApp()

		e.WaitForStop()
	}()

	front.Run()
}