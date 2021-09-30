//
//  ImageRendererDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 24.08.2021.
//

import Foundation
import MetalKit

struct ImageVertexDescriptor {
    var position: SIMD3<Float>
    var texCoord: SIMD2<Float>
}

class ImageRendererDelegate: BasePrimitiveRendererDelegate {
    private var textureLoader: MTKTextureLoader?
    
    override func onStart() {
        let imageVertexFunction = library.makeFunction(name: "image_vertex")!
        let imageFragmentFunction = library.makeFunction(name: "image_fragment")!
        
        let imagePipelineDescriptor = makeDefaultRenderPipelineDescriptor(imageVertexFunction, imageFragmentFunction)
        
        let imageVertexDescriptor = MTLVertexDescriptor()
        imageVertexDescriptor.attributes[0].format = .float3
        imageVertexDescriptor.attributes[0].offset = 0
        imageVertexDescriptor.attributes[0].bufferIndex = 0
        imageVertexDescriptor.attributes[1].format = .float2
        imageVertexDescriptor.attributes[1].offset = MemoryLayout<SIMD3<Float>>.stride
        imageVertexDescriptor.attributes[1].bufferIndex = 0
        
        imageVertexDescriptor.layouts[0].stride = MemoryLayout<ImageVertexDescriptor>.stride
        
        imagePipelineDescriptor.vertexDescriptor = imageVertexDescriptor

        pipelineState = try! device.makeRenderPipelineState(descriptor: imagePipelineDescriptor)
        
        textureLoader = MTKTextureLoader(device: device)
    }

    override func onRender(_ ctx: CliPrimitiveRenderingContext?) {
        guard let context = ctx,
              var primitive = PrimitivesRegistry.shared.getPrimitive(byId: context.primitiveId),
              let image = context.imagePrimitiveData,
              let tlp = image.tlPosition,
              let brp = image.brPosition else {return}
        
        if context.redraw {
            if primitive.textures == nil {
                primitive.textures = []
                
                for i in 0..<image.getBitmapCount() {
                    guard let bitmap = image.bitmap(at: i),
                          let textureData = bitmap.getPixels() else {return}
                        
                    let textureDescriptor = MTLTextureDescriptor()
                    textureDescriptor.pixelFormat = .rgba8Uint
                    textureDescriptor.width = bitmap.getWidth()
                    textureDescriptor.height = bitmap.getHeight()

                    let texture = device.makeTexture(descriptor: textureDescriptor)!
                    let texturePointer = textureData.withUnsafeBytes { (unsafePointer: UnsafeRawBufferPointer) -> UnsafeRawPointer in
                        return unsafePointer.baseAddress!
                    }

                    let region = MTLRegion(origin: MTLOrigin(x: 0, y: 0, z: 0), size: MTLSize(width: bitmap.getWidth(), height: bitmap.getHeight(), depth: 1))
                    texture.replace(region: region, mipmapLevel: 0, withBytes: texturePointer, bytesPerRow: bitmap.getWidth() * 4)
                    
                    primitive.textures!.append(texture)
                }
                
//                print("Generated textures \(primitive.textures!.count)")
            }
            
            let verticies: [ImageVertexDescriptor] = [
                ImageVertexDescriptor(position: SIMD3<Float>(tlp.x, tlp.y, tlp.z), texCoord: SIMD2<Float>(0, 0)),
                ImageVertexDescriptor(position: SIMD3<Float>(tlp.x, brp.y, tlp.z), texCoord: SIMD2<Float>(0, 1)),
                ImageVertexDescriptor(position: SIMD3<Float>(brp.x, brp.y, tlp.z), texCoord: SIMD2<Float>(1, 1)),
                ImageVertexDescriptor(position: SIMD3<Float>(brp.x, tlp.y, tlp.z), texCoord: SIMD2<Float>(1, 0)),
            ]
            
            let indicies: [UInt16] = [
                0, 1, 2,
                2, 3, 0,
            ]
            
            primitive.vertexBuffer = device.makeBuffer(bytes: verticies, length: verticies.count * MemoryLayout<ImageVertexDescriptor>.stride, options: [])
            primitive.indexBuffer = device.makeBuffer(bytes: indicies, length: indicies.count * MemoryLayout<UInt16>.stride, options: [])
            
            let samplerDescriptor = MTLSamplerDescriptor()
            samplerDescriptor.minFilter = .linear
            samplerDescriptor.magFilter = .linear
            primitive.sampler = device.makeSamplerState(descriptor: samplerDescriptor)
            
            PrimitivesRegistry.shared.setPrimitive(byId: context.primitiveId, withData: primitive)
        }
        
        var uniform = UniformBuffer(projection: Projection.current)
        
//        print("Rendering image. Texture: \(image.index)")
        
        RendererDelegate.renderEncoder.setRenderPipelineState(pipelineState)
        RendererDelegate.renderEncoder.setVertexBuffer(primitive.vertexBuffer!, offset: 0, index: 0)
        RendererDelegate.renderEncoder.setVertexBytes(&uniform, length: UniformBuffer.stride, index: 1)
        RendererDelegate.renderEncoder.setFragmentTexture(primitive.textures![image.index], index: 0)
        RendererDelegate.renderEncoder.setFragmentSamplerState(primitive.sampler!, index: 0)
        RendererDelegate.renderEncoder.drawIndexedPrimitives(type: .triangle, indexCount: 6, indexType: .uint16, indexBuffer: primitive.indexBuffer!, indexBufferOffset: 0)
    }
}
