//
//  TextRendererDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 30.09.2021.
//

import Foundation
import MetalKit
import Amphion

struct TextVertexDescriptor {
    var position: SIMD3<Float>
    var texCoord: SIMD2<Float>
    var color: SIMD4<Float>
}

class TextRendererDelegate: BasePrimitiveRendererDelegate {
    
    let textures: GlyphTexturesRegistry
    
    override init(_ lib: MTLLibrary, _ dev: MTLDevice) {
        textures = GlyphTexturesRegistry(device: dev)
        super.init(lib, dev)
    }
    
    override func onStart() {
        let textVertexFunction = library.makeFunction(name: "text_vertex")!
        let textFragemntFunction = library.makeFunction(name: "text_fragment")!
        
        let textPipelineDescriptor = makeDefaultRenderPipelineDescriptor(textVertexFunction, textFragemntFunction)
        textPipelineDescriptor.colorAttachments[0].pixelFormat = .bgra8Unorm
        textPipelineDescriptor.colorAttachments[0].isBlendingEnabled = true
        textPipelineDescriptor.colorAttachments[0].alphaBlendOperation = .add
        textPipelineDescriptor.colorAttachments[0].rgbBlendOperation = .add
        textPipelineDescriptor.colorAttachments[0].sourceRGBBlendFactor = .sourceAlpha
        textPipelineDescriptor.colorAttachments[0].sourceAlphaBlendFactor = .sourceAlpha
        textPipelineDescriptor.colorAttachments[0].destinationRGBBlendFactor = .oneMinusSourceColor
        textPipelineDescriptor.colorAttachments[0].destinationAlphaBlendFactor = .oneMinusSourceAlpha
        
        let textVertexDescriptor = MTLVertexDescriptor()
        textVertexDescriptor.attributes[0].format = .float3
        textVertexDescriptor.attributes[0].offset = 0
        textVertexDescriptor.attributes[0].bufferIndex = 0
        
        textVertexDescriptor.attributes[1].format = .float2
        textVertexDescriptor.attributes[1].offset = MemoryLayout<SIMD3<Float>>.stride
        textVertexDescriptor.attributes[1].bufferIndex = 0
        
        textVertexDescriptor.attributes[2].format = .float4
        textVertexDescriptor.attributes[2].offset = MemoryLayout<SIMD3<Float>>.stride + MemoryLayout<SIMD2<Float>>.stride
        textVertexDescriptor.attributes[2].bufferIndex = 0

        textVertexDescriptor.layouts[0].stride = MemoryLayout<TextVertexDescriptor>.stride

        textPipelineDescriptor.vertexDescriptor = textVertexDescriptor

        pipelineState = try! device.makeRenderPipelineState(descriptor: textPipelineDescriptor)
    }
    
    override func onSetPrimitive(_ ctx: CliPrimitiveRenderingContext?) {
        super.onSetPrimitive(ctx)
        print("Setting text primitive: \(String(describing: ctx?.textPrimitiveData?.text))")
    }
    
    override func onRender(_ ctx: CliPrimitiveRenderingContext?) {
        guard let context = ctx,
              let primitive = PrimitivesRegistry.shared.getPrimitive(byId: context.primitiveId),
              let text = context.textPrimitiveData,
              let tlPosition = text.tlPosition,
              let atext = text.provider?.getAText(),
              let color = text.textColor else {return}
        
        if context.redraw {
            primitive.vertexBuffers.removeAll()
            primitive.indexBuffers.removeAll()
            
            if primitive.sampler == nil {
                let samplerDescriptor = MTLSamplerDescriptor()
                samplerDescriptor.minFilter = .linear
                samplerDescriptor.magFilter = .linear
                primitive.sampler = device.makeSamplerState(descriptor: samplerDescriptor)
            }
        }
        
//        print("Rendering text \(text.text)")
        
        for i in 0..<atext.getCharsCount() {
            guard let char = atext.getCharAt(i) else {continue}
            if !char.isVisible() {
                continue
            }
            let glyph = char.getGlyph()!
            
            if context.redraw {
                let tlp = CliNewVector3(Float(char.getX()), Float(char.getY()), 0)!
                let brp = CliNewVector3(tlp.x + Float(glyph.getWidth()), tlp.y + Float(glyph.getHeight()), 0)!
                
                let simdColor = SIMD4<Float>(color.x, color.y, color.z, color.w)
                
                print("Color: {\(color.x), \(color.y), \(color.z), \(color.w)} simd: {\(simdColor.x), \(simdColor.y), \(simdColor.z), \(simdColor.w)}")
                
                let vertices: [TextVertexDescriptor] = [
                    TextVertexDescriptor(position: SIMD3<Float>(tlp.x, tlp.y, 0), texCoord: SIMD2<Float>(0, 0), color: simdColor),
                    TextVertexDescriptor(position: SIMD3<Float>(tlp.x, brp.y, 0), texCoord: SIMD2<Float>(0, 1), color: simdColor),
                    TextVertexDescriptor(position: SIMD3<Float>(brp.x, brp.y, 0), texCoord: SIMD2<Float>(1, 1), color: simdColor),
                    TextVertexDescriptor(position: SIMD3<Float>(brp.x, tlp.y, 0), texCoord: SIMD2<Float>(1, 0), color: simdColor),
                ]
                
                let indicies: [UInt16] = [
                    0, 1, 2,
                    2, 3, 0,
                ]
                
                primitive.vertexBuffers.append(device.makeBuffer(bytes: vertices, length: vertices.count * MemoryLayout<TextVertexDescriptor>.stride, options: [])!)
                primitive.indexBuffers.append(device.makeBuffer(bytes: indicies, length: indicies.count * MemoryLayout<UInt16>.size, options: [])!)
                
                print("redraw text \(text.text)")
            }
            
            PrimitivesRegistry.shared.setPrimitive(byId: context.primitiveId, withData: primitive)
            
            let texture = textures.ensureGlyphTextureIsGenerated(forGlyph: glyph)
            
            var uniform = UniformBuffer(projection: Projection.current)
            
//            print("Rendering char \(char.getRune())")
            
            RendererDelegate.renderEncoder.setRenderPipelineState(pipelineState)
            RendererDelegate.renderEncoder.setVertexBuffer(primitive.vertexBuffers[i], offset: 0, index: 0)
            RendererDelegate.renderEncoder.setVertexBytes(&uniform, length: UniformBuffer.stride, index: 1)
            RendererDelegate.renderEncoder.setFragmentTexture(texture, index: 0)
            RendererDelegate.renderEncoder.setFragmentSamplerState(primitive.sampler!, index: 0)
            RendererDelegate.renderEncoder.drawIndexedPrimitives(type: .triangle, indexCount: 6, indexType: .uint16, indexBuffer: primitive.indexBuffers[i], indexBufferOffset: 0)
        }
    }
}
