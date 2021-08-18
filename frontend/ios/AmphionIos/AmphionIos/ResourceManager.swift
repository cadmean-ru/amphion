//
//  ResourceManager.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation

class ResourceManager : NSObject, CliResourceManagerDelegateProtocol {
    func readFile(_ path: String?) throws -> Data {
        guard let path2 = path,
              let url = URL(string: path2) else { return Data() }
        
        let fullPath = Bundle.main.url(forResource: url.deletingPathExtension().lastPathComponent, withExtension: url.pathExtension, subdirectory: url.deletingLastPathComponent().path)!
        return try Data(contentsOf: fullPath)
    }
}
