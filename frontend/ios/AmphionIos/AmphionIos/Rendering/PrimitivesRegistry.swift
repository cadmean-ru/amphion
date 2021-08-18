//
//  PrimitivesRegistry.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation

class PrimitivesRegistry {
    static let shared = PrimitivesRegistry()
    
    private var primitives: [Int:MetalPrimitive] = [:]
    
    func getPrimitive(byId id: Int) -> MetalPrimitive? {
        return primitives[id]
    }
    
    func setPrimitive(byId id: Int, withData p: MetalPrimitive) {
        primitives[id] = p
    }
    
    func createPrimitive(withId id: Int) {
        primitives[id] = MetalPrimitive()
    }
    
    func primitiveExists(_ id: Int) -> Bool {
        return getPrimitive(byId: id) != nil
    }
    
    func removePrimitive(withId id: Int) {
        primitives.removeValue(forKey: id)
    }
}
