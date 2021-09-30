//
//  AmphionGestureRecognizer.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import UIKit

class AmphionGestureRecognizer: UIGestureRecognizer {
    override func touchesBegan(_ touches: Set<UITouch>, with event: UIEvent) {
        guard let touch = touches.first else { return }
        let position = touch.location(in: view)
        FrontendDelegate.shared.sendCallback(of: -108, withStringData: "\(Int(position.x));\(Int(position.y));0")
    }
    
    override func touchesMoved(_ touches: Set<UITouch>, with event: UIEvent) {
        guard let touch = touches.first else { return }
        let position = touch.location(in: view)
        FrontendDelegate.shared.sendCallback(of: -110, withStringData: "\(Int(position.x));\(Int(position.y))")
    }
    
    override func touchesEnded(_ touches: Set<UITouch>, with event: UIEvent) {
        guard let touch = touches.first else { return }
        let position = touch.location(in: view)
        FrontendDelegate.shared.sendCallback(of: -109, withStringData: "\(Int(position.x));\(Int(position.y));0")
    }
}
