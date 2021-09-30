//
//  UniformBuffer.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 01.09.2021.
//

import Foundation
import MetalKit

struct UniformBuffer {
    let projection: matrix_float4x4
}

extension UniformBuffer {
    static var stride: Int {
        get {
            return MemoryLayout<UniformBuffer>.stride
        }
    }
}
