//
//  ProjectionTest.swift
//  AmphionIosTests
//
//  Created by Алексей Крицков on 30.09.2021.
//

import XCTest
@testable import AmphionIos

class ProjectionTest: XCTestCase {
    
    override func setUpWithError() throws {
        Projection.calculate(width: 400, height: 800)
    }
    
    func testApply() {
        let v = SIMD4<Float>(0, 0, 0, 1)
        let r = Projection.apply(to: v)
        print(r)
    }
}
