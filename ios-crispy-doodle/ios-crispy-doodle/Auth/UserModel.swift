//
//  UserModel.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation
import Observation

struct User: Codable, Identifiable, Equatable {
    var id: String
    var name: String
    var email: String
    var password: String?
    var online: Bool
    var channels: [String] = []
    var created: Double
    var updated: Double

    func encodeUser(_ user: User) -> Data? {
        let encoder = JSONEncoder()
        do {
            return try encoder.encode(user)
        } catch {
            print("Encoding error: \(error)")
            return nil
        }
    }

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


//ios1
//
//
