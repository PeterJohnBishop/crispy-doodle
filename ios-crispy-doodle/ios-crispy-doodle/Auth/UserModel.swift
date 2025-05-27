//
//  UserModel.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation

struct User: Codable, Identifiable, Equatable {
    var id: String
    var name: String
    var email: String
    var password: String?
    var online: Bool
    var channels: [String] = []
    var created: Double
    var updated: Double
}

extension User {
    static func decodeUser(from data: Data) throws -> User {
        let decoder = JSONDecoder()
        return try decoder.decode(User.self, from: data)
    }
}

extension User {
    static  func encodeUser(_ user: User) throws -> Data?{
        let encoder = JSONEncoder()
        return try encoder.encode(user)
    }
}


