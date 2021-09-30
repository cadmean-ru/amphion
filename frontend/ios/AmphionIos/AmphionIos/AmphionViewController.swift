//
//  AmphionViewController.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import UIKit

open class AmphionViewController: UIViewController {
    open override var preferredStatusBarStyle: UIStatusBarStyle {
        return .darkContent
    }
    
    open override func viewDidLoad() {
        super.viewDidLoad()
        
        setNeedsStatusBarAppearanceUpdate()
        
        FrontendDelegate.shared.setView(view)
        IosCliAmphionInit(FrontendDelegate.shared, ResourceManager(), RendererDelegate(for: view))
        FrontendDelegate.shared.sendCallback(of: -111, withStringData: "")
        
        view.addGestureRecognizer(AmphionGestureRecognizer())
    }
}
