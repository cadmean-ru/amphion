//
//  ViewExtension.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 30.09.2021.
//

import Foundation
import UIKit
import Amphion

extension UIView {
    func getWindowSize() -> CliVector3 {
        let v = CliNewVector3(0, 0, 0)!
        
        if (self.window != nil) {
            let windowFrame = self.window!.frame
            v.x = Float(windowFrame.size.width)
            v.y = Float(windowFrame.size.height)
        } else {
            let screenBounds = UIScreen.main.bounds
            v.x = Float(screenBounds.size.width)
            v.y = Float(screenBounds.size.height)
        }
        
        return v
    }
}
