//
//  FrontendDelegate.swift
//  AmphionIos
//
//  Created by Алексей Крицков on 18.08.2021.
//

import Foundation
import UIKit

class FrontendDelegate : NSObject, CliFrontendDelegateProtocol, DispatchWorkDispatcherProtocol {
    static let shared = FrontendDelegate()
    
    private var view: UIView?
    private var callbackDispatcher: DispatchMessageDispatcherProtocol?
    
    func setView(_ view: UIView) {
        self.view = view
    }
    
    func init_() {
        print("Initializing iOS frontend")
    }
    
    func run() {
        print("Running iOS frontend")
    }
    
    func commencePanic(_ reason: String?, msg: String?) {
        print("----- iOS Frontend Panic -----")
        print("Reason: \(String(describing: reason))")
        print(String(describing: msg))
    }
    
    func getAppData() -> Data? {
        let filePath = Bundle.main.url(forResource: "app", withExtension: "yaml")!
        return try? Data(contentsOf: filePath)
    }
    
    func getContext() -> CliContext? {
        let ctx = CliContext()
        ctx.screenSize = view?.getWindowSize()
        return ctx
    }
    
    func getMainThreadDispatcher() -> DispatchWorkDispatcherProtocol? {
        return self
    }
    
    func getRenderingThreadDispatcher() -> DispatchWorkDispatcherProtocol? {
        return self
    }
    
    
    func setCallbackDispatcher(_ dispatcher: DispatchMessageDispatcherProtocol?) {
        callbackDispatcher = dispatcher
    }
    
    func execute(_ item: DispatchWorkItemProtocol?) {
        DispatchQueue.main.async {
            item?.execute()
        }
    }
    
    func sendCallback(of code: Int, withStringData data: String) {
        callbackDispatcher?.send(DispatchNewMessageWithStringData(code, data))
    }
}
