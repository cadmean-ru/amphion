//
//  DeviceExtension.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 30.09.2021.
//

import Foundation
import MetalKit

extension MTLDevice {
    func makeTexture(descriptor textureDescriptor: MTLTextureDescriptor, pixels: Data, width: Int, height: Int, bytesPerPixel: Int) -> MTLTexture {
        let texture = self.makeTexture(descriptor: textureDescriptor)!
        
        let texturePointer = pixels.withUnsafeBytes { (unsafePointer: UnsafeRawBufferPointer) -> UnsafeRawPointer in unsafePointer.baseAddress!}
        
        let region = MTLRegion(origin: MTLOrigin(x: 0, y: 0, z: 0), size: MTLSize(width: width, height: height, depth: 1))
        texture.replace(region: region, mipmapLevel: 0, withBytes: texturePointer, bytesPerRow: bytesPerPixel * width)
        
        return texture
    }
}
