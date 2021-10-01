//
//  GlyphTexturesRegistry.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 30.09.2021.
//

import Foundation
import MetalKit
import Amphion

class GlyphTextureHolder {
    let glyph: AtextGlyph
    let texture: MTLTexture
    
    init(holding glyph: AtextGlyph, withTexture texture: MTLTexture) {
        self.glyph = glyph
        self.texture = texture
    }
    
    var fontName: String {
        get {
            return glyph.getFace()!.getFont()!.getName()
        }
    }
    
    var fontSize: Int {
        get {
            return glyph.getFace()!.getSize()
        }
    }
    
    var char: Int32 {
        get {
            return glyph.getRune()
        }
    }
    
    func matches(char: Int32, ofSize size: Int, inFont font: String) -> Bool {
        return fontName == font && fontSize == size && self.char == char
    }
}

class GlyphTexturesRegistry {
    var glyphs: [GlyphTextureHolder] = []
    let device: MTLDevice
    
    init(device: MTLDevice) {
        self.device = device
    }
    
    func ensureGlyphTextureIsGenerated(forGlyph glyph: AtextGlyph) -> MTLTexture {
        let t = getGlyphTexture(forChar: glyph.getRune(), ofSize: glyph.getFace()!.getSize(), inFont: glyph.getFace()!.getFont()!.getName())
        if t != nil {
            return t!
        }
        
        let textureDescriptor = MTLTextureDescriptor()
        textureDescriptor.pixelFormat = .r8Uint
        textureDescriptor.width = glyph.getWidth()
        textureDescriptor.height = glyph.getHeight()
        
        let texture = device.makeTexture(descriptor: textureDescriptor, pixels: glyph.getPixels()!, width: glyph.getWidth(), height: glyph.getHeight(), bytesPerPixel: 1)
        
        glyphs.append(GlyphTextureHolder(holding: glyph, withTexture: texture))
        
        return texture
    }
    
    func getGlyphTexture(forChar char: Int32, ofSize size: Int, inFont font: String) -> MTLTexture? {
        for g in glyphs {
            if g.matches(char: char, ofSize: size, inFont: font) {
                return g.texture
            }
        }
        
        return nil
    }
}
