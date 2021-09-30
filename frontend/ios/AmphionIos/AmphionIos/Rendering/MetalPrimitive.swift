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

struct MetalPrimitive {
    var data: Any?
    var redraw: Bool
    var z: Float
    var kind: PrimitiveKind
    var vertexBuffer: MTLBuffer?
    var indexBuffer: MTLBuffer?
    var textures: [MTLTexture]?
    var sampler: MTLSamplerState?
    var uniformBuffer: MTLBuffer?
    
    init() {
        data = nil
        redraw = false
        z = 0
        kind = .Unset
    }
}
