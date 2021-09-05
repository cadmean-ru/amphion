//
//  RendererDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import UIKit
import MetalKit


class RendererDelegate : NSObject, CliRendererDelegateProtocol {
    
    private var view: UIView!
    
    internal var device: MTLDevice!
    private var metalLayer: CAMetalLayer!

    private var commandQueue: MTLCommandQueue!
    internal static var renderEncoder: MTLRenderCommandEncoder!
    
    private var drawable: CAMetalDrawable?
    private var commandBuffer: MTLCommandBuffer?
    
    init(for view: UIView!) {
        self.view = view
    }
    
    func onPrepare() {
        device = MTLCreateSystemDefaultDevice()
        
        metalLayer = CAMetalLayer()
        metalLayer.device = device
        metalLayer.pixelFormat = .bgra8Unorm
        metalLayer.framebufferOnly = true
        metalLayer.frame = view.layer.frame
        view.layer.addSublayer(metalLayer)
     
        //Default libaray and pipeline
        let amphionBundle = Bundle(for: RendererDelegate.self)
        guard let library = try? device.makeDefaultLibrary(bundle: amphionBundle) else {
            print("Could not load default library")
            return
        }
        
        commandQueue = device.makeCommandQueue()
        
        IosCliRegisterPrimitiveRendererDelegate(CliPrimitiveRectangle, RectRendererDelegate(library, device))
        IosCliRegisterPrimitiveRendererDelegate(CliPrimitiveImage, ImageRendererDelegate(library, device))
//        IosCliRegisterPrimitiveRendererDelegate(CliPrimitiveText, TextRendererDelegate(library, device))
    }
    
    func onPerformRenderingStart() {
        guard let drawable = metalLayer!.nextDrawable() else { return }
        self.drawable = drawable
        let renderPassDescriptor = MTLRenderPassDescriptor()
        renderPassDescriptor.colorAttachments[0].texture = drawable.texture
        renderPassDescriptor.colorAttachments[0].loadAction = .clear
        renderPassDescriptor.colorAttachments[0].clearColor = MTLClearColor(red: 1, green: 1, blue: 1, alpha: 1)
        
        commandBuffer = commandQueue.makeCommandBuffer()
        RendererDelegate.renderEncoder = commandBuffer?.makeRenderCommandEncoder(descriptor: renderPassDescriptor)!
        
        print("On perform rendering start ios")
    }
    
    func onPerformRenderingEnd() {
        RendererDelegate.renderEncoder.endEncoding()
        
        commandBuffer?.present(drawable!)
        commandBuffer?.commit()
        
        print("On perform rendering end ios")
    }
    
    func onClear() {
        
    }
    
    func onStop() {
        
    }
}
