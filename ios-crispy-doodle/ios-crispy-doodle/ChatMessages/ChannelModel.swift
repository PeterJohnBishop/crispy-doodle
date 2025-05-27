//
//  ChannelModel.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/26/25.
//

import Foundation

struct Channel: Codable, Identifiable, Equatable {
    var id: String
    var title: String
    var messages: [String] = []
    var created: Double
    var updated: Double
}

extension Channel {
    static func decodeChannel(from data: Data) throws -> Channel {
        let decoder = JSONDecoder()
        return try decoder.decode(Channel.self, from: data)
    }
}

extension Channel {
    static  func encodeChannel(_ channel: Channel) throws -> Data?{
        let encoder = JSONEncoder()
        return try encoder.encode(channel)
    }
}

