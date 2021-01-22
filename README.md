# Amphion engine

Amphion is a component oriented application engine.
It's goal is to simplify multiplatform apps development, providing easy single-language (Go) solution, 
adapting component oriented programming techniques.

**This is still a prototype/technology preview.**

## How it works

An amphion based application consists of code, resources and some config files.

The engine consists of the kernel and a frontend, that is specific for different platforms.

The frontend is used for:
* rendering;
* getting user input and handling it to the kernel.

When the engine is started, it tries to load application config file from well-known location (depends on frontend).
Based on the config, the engine then displays the main scene.

### Scene

An application consists of scenes. Each scene containes multiplse scene objects.
A **scene object** is a basic building block of a scene. It has some size and position in the scene, but it does nothing.
The scene itself is also a **scene object**.

Scene can be serialized, saved and than loaded as a resource.

Only one scene can be showed to the user at a time.

### Components

A **component** is a piece of functionality that can be attached to a **scene object**.

For example, to display a rect in the scene, you need to create a scene object, specify it's position and size 
and than add a component that draws a rectangle on the screen.
A component that draws something on the screen is called **view**.

A component is ultimately a Go struct, that implements *engine.Component* interface.

Component interface requires the following lifecycle methods:
* OnInit(ctx engine.InitContext)
* OnStart()
* OnStop()

It also requires GetName() string method, that returns the name of the component.

Lifecycle method are called by the kernel at specific stages of component's lifecycle.

Component can also implment additional interfaces, introducing additional lifecycle methods: 
* UpdatingComponent{*OnUpdate(ctx engine.UpdateContext)*}, 
* ViewComponent{*OnDraw(ctx engine.DrawingContext)*; *ForceRedraw()*}.

### Scene lifecycle

When the scene is loaded, the kernel starts the update routine.
The update routine is basically an infinite loop, where the scene's logic is executed.

Each iteration of the loop consists of the following stages:
* initializing components, that hasn't been components yet (OnInit);
* starting components if needed (OnStart);
* stopping components if needed (OnStop);
* updating components if requested (OnUpdate);
* drawing components if requested (OnDraw);
* perfoming actual rendering.

The loop won't be executed every frame.
It will be executed once when the scene is shown. 
Than it will wait for one of the components to explicitly request next update.

OnInit is called only once for each component after it was created.
OnStart is called every time the component is enabled and after it was created.
OnStop is called every time the component is disabled.
OnUpdate is called on every update routine iteration.
OnDraw is called on update routine iteration if rendering was requested.

### Rendering

TODO

## Supported platforms

|Platform|Status         |
|--------|---------------|
|Web     |development    |
|Windows |development    |
|Linux   |development    |
|macOS   |development*   |
|Android |planned        |
|iOS     |planned        |

*Currently, supported through OpenGL, Metal support planned.

## Features

|Feature|Status|
|-------|------|
|Basic views|In development|
|Navigation|planned|
|Animations|planned|
|RPC|In development|
|Native API|planned|
|Rendering API|planned|

## Contributing

Feel free to submit issues or create pull requests following fork-and-pull git workflow.
