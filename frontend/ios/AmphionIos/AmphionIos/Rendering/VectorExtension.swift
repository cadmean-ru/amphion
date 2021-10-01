//
//  VectorExtension.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 01.10.2021.
//

import Foundation
import MetalKit

extension CliVector4 {
    var simd: SIMD4<Float> {
        get {
            SIMD4<Float>(self.x, self.y, self.z, self.w)
        }
    }
}

extension CliVector3 {
    var simd: SIMD3<Float> {
        get {
            SIMD3<Float>(self.x, self.y, self.z)
        }
    }
}
