//
//  Item.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
