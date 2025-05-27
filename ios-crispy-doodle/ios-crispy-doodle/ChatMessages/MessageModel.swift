//
//  MessageModel.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/26/25.
//

import Foundation

struct Message: Codable, Identifiable, Equatable {
    var id: String
    var sender: String
    var text: String
    var images: [String] = []
    var created: Double
    var updated: Double
}

extension Message {
    static func decodeMessage(from data: Data) throws -> Message {
        let decoder = JSONDecoder()
        return try decoder.decode(Message.self, from: data)
    }
}

extension Message {
    static  func encodeMessage(_ message: Message) throws -> Data?{
        let encoder = JSONEncoder()
        return try encoder.encode(message)
    }
}
