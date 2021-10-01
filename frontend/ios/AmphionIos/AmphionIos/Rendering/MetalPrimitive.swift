//
//  MetalPrimitiveState.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import MetalKit

enum PrimitiveKind {
    case Rectangle
    case Ellipse
    case Triangle
    case Text
    case Image
    case Unset
}

class MetalPrimitive {
    var data: Any?
    var kind: PrimitiveKind
    var vertexBuffers: [MTLBuffer] = []
    var indexBuffers: [MTLBuffer] = []
    var textures: [MTLTexture]?
    var sampler: MTLSamplerState?
    
    init() {
        data = nil
        kind = .Unset
    }
    
    var vertexBuffer: MTLBuffer? {
        get {
            vertexBuffers.first
        }
        set {
            if vertexBuffers.isEmpty {
                vertexBuffers.append(newValue!)
            } else {
                vertexBuffers[0] = newValue!
            }
        }
    }
    
    var indexBuffer: MTLBuffer? {
        get {
            indexBuffers.first
        }
        set {
            if indexBuffers.isEmpty {
                indexBuffers.append(newValue!)
            } else {
                indexBuffers[0] = newValue!
            }
        }
    }
}
