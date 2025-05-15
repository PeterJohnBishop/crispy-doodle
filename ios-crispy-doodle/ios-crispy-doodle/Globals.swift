//
//  Globals.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation

struct Global {
    
    static var baseURL: String = "http://localhost:8080"

}

extension NSError: LocalizedError {
    
    public var errorDescription: String? { localizedDescription }
    
    public var failureReason: String? { localizedFailureReason }
    
    public var recoverySuggestion: String? { localizedRecoverySuggestion }
    
    // public var helpAnchor: String? { get }  // Goes as is.
}
