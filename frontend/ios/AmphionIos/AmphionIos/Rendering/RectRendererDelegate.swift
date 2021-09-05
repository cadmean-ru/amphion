//
//  RectRendererDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import MetalKit

struct RectVertexDescriptor {
    var position: SIMD3<Float>
    var color: SIMD4<Float>
}

class RectRendererDelegate : BasePrimitiveRendererDelegate {
    override func onStart() {
        let rectVertexFunction = library.makeFunction(name: "rect_vertex")!
        let rectFragmentFunction = library.makeFunction(name: "rect_fragment")!
        
        let rectPipelineDescriptor = makeDefaultRenderPipelineDescriptor(rectVertexFunction, rectFragmentFunction)
        
        let rectVertexDescriptor = MTLVertexDescriptor()
        rectVertexDescriptor.attributes[0].format = .float3
        rectVertexDescriptor.attributes[0].offset = 0
        rectVertexDescriptor.attributes[0].bufferIndex = 0
        rectVertexDescriptor.attributes[1].format = .float4
        rectVertexDescriptor.attributes[1].offset = MemoryLayout<SIMD3<Float>>.stride
        rectVertexDescriptor.attributes[1].bufferIndex = 0
        
        rectVertexDescriptor.layouts[0].stride = MemoryLayout<RectVertexDescriptor>.stride
        
        rectPipelineDescriptor.vertexDescriptor = rectVertexDescriptor

        pipelineState = try! device.makeRenderPipelineState(descriptor: rectPipelineDescriptor)
    }
    
    override func onRender(_ ctx: CliPrimitiveRenderingContext?) {
        guard let context = ctx,
              var primitive = PrimitivesRegistry.shared.getPrimitive(byId: context.primitiveId),
              let geometry = context.geometryPrimitiveData,
              let tlp = geometry.tlPositionN,
              let brp = geometry.brPositionN,
              let fillColor = geometry.fillColorN else {return}
        
        if (context.redraw) {
            let simdColor = SIMD4<Float>(fillColor.x, fillColor.y, fillColor.z, fillColor.w)
            let verticies: [RectVertexDescriptor] = [
                RectVertexDescriptor(position: SIMD3<Float>(tlp.x, tlp.y, tlp.z), color: simdColor),
                RectVertexDescriptor(position: SIMD3<Float>(tlp.x, brp.y, tlp.z), color: simdColor),
                RectVertexDescriptor(position: SIMD3<Float>(brp.x, brp.y, tlp.z), color: simdColor),
                RectVertexDescriptor(position: SIMD3<Float>(brp.x, tlp.y, tlp.z), color: simdColor),
            ]
            
            let indicies: [UInt16] = [
                0, 1, 2,
                2, 3, 0,
            ]
            
            let 🍆: [Float] = [
                
            ]
            
            primitive.vertexBuffer = device.makeBuffer(bytes: verticies, length: verticies.count * MemoryLayout<RectVertexDescriptor>.stride, options: [])
            primitive.indexBuffer = device.makeBuffer(bytes: indicies, length: indicies.count * MemoryLayout<UInt16>.stride, options: [])
            
            primitive.uniformBuffer = device.makeBuffer(bytes: 🍆, length: MemoryLayout<UniformBuffer>.size, options: [])
            
            PrimitivesRegistry.shared.setPrimitive(byId: context.primitiveId, withData: primitive)
        }
        
        RendererDelegate.renderEncoder.setRenderPipelineState(pipelineState)
        RendererDelegate.renderEncoder.setVertexBuffer(primitive.vertexBuffer!, offset: 0, index: 0)
        RendererDelegate.renderEncoder.drawIndexedPrimitives(type: .triangle, indexCount: 6, indexType: .uint16, indexBuffer: primitive.indexBuffer!, indexBufferOffset: 0)
    }
}
