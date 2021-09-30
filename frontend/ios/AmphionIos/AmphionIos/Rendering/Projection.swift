//
//  Projection.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 29.09.2021.
//

import Foundation
import UIKit
import MetalKit

class Projection {
    static var current: matrix_float4x4 = .init()
    
    static func calculate(for view: UIView) {
        let wSize = view.getWindowSize()
        calculate(width: wSize.x, height: wSize.y)
    }
    
    static func calculate(width: Float, height: Float) {
        let xs = width
        let ys = height
        let zs = 2
        let c1 = 2.0 / Float(xs)
        let c2 = -2.0 / Float(ys)
        let c3 = -2.0 / Float(zs)
        
//        Projection.current = [
//            c1,  0,  0, -1,
//             0, c2,  0,  1,
//             0,  0, c3,  0,
//             0,  0,  0,  1
//        ]
        
        Projection.current = .init([
            SIMD4<Float>(c1, 0, 0, 0),
            SIMD4<Float>(0, c2, 0, 0),
            SIMD4<Float>(0, 0, c3, 0),
            SIMD4<Float>(-1, 1, 0, 1),
        ])
    }
    
    static func apply(to vector: SIMD4<Float>) -> SIMD4<Float> {
        return simd_mul(Projection.current, vector)
    }
    
    static func apply(to vector: SIMD3<Float>) -> SIMD3<Float> {
        let v4 = simd_mul(Projection.current, SIMD4<Float>(vector, 1))
        return SIMD3<Float>(v4.x, v4.y, v4.z)
    }
}
