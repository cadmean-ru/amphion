//
//  BasePrimitiveRendererDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import MetalKit

class BasePrimitiveRendererDelegate : NSObject, CliPrimitiveRendererDelegateProtocol {
    var pipelineState: MTLRenderPipelineState!
    var library: MTLLibrary
    var device: MTLDevice
    
    init(_ lib: MTLLibrary, _ dev: MTLDevice) {
        library = lib
        device = dev
    }
    
    func onStart() {
        
    }
    
    func onSetPrimitive(_ ctx: CliPrimitiveRenderingContext?) {
        guard let c = ctx else { return }
        if !PrimitivesRegistry.shared.primitiveExists(c.primitiveId) {
            PrimitivesRegistry.shared.createPrimitive(withId: c.primitiveId)
        }
        
//        print("On set primitive \(String(describing: ctx?.primitiveId))")
    }
    
    func onRemovePrimitive(_ ctx: CliPrimitiveRenderingContext?) {
        guard let context = ctx else { return }
        PrimitivesRegistry.shared.removePrimitive(withId: context.primitiveId)
    }
    
    func onRender(_ ctx: CliPrimitiveRenderingContext?) {
        
    }
    
    func onStop() {
        
    }
    
    func makeDefaultRenderPipelineDescriptor(_ vertexFunxtion: MTLFunction, _ fragmentFunction: MTLFunction) -> MTLRenderPipelineDescriptor {
        let descriptor = MTLRenderPipelineDescriptor()
        
        descriptor.vertexFunction = vertexFunxtion
        descriptor.fragmentFunction = fragmentFunction
        descriptor.colorAttachments[0].pixelFormat = .bgra8Unorm
        descriptor.colorAttachments[0].isBlendingEnabled = true
        descriptor.colorAttachments[0].alphaBlendOperation = .add
        descriptor.colorAttachments[0].rgbBlendOperation = .add
        descriptor.colorAttachments[0].sourceRGBBlendFactor = .sourceAlpha
        descriptor.colorAttachments[0].sourceAlphaBlendFactor = .sourceAlpha
        descriptor.colorAttachments[0].destinationRGBBlendFactor = .oneMinusSourceAlpha
        descriptor.colorAttachments[0].destinationAlphaBlendFactor = .oneMinusSourceAlpha
        
        return descriptor
    }
}
