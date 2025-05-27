//
//  UserViewModel.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation
import Observation

@MainActor
@Observable class UserViewModel: ObservableObject {
    var user: User = User(id: "", name: "", email: "", online: false, created: 0, updated: 0)
    var users: [User] = []
    var error: String?
    var isLoading: Bool = false
    
    func RegisterUser() async -> Bool {
        guard let url = URL(string: "\(Global.baseURL)/register") else {
            print("Invalid URL: \(Global.baseURL)/register")
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "name": user.name,
            "email": user.email,
            "password": user.password ?? ""
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else {
            print("Failed to serialize request body")
            return false
        }

        request.httpBody = jsonData
        print("Sending request to: \(url.absoluteString)")
        print("Request body: \(String(data: jsonData, encoding: .utf8) ?? "Invalid JSON")")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse {
                print("Received response with status code: \(httpResponse.statusCode)")

                if httpResponse.statusCode == 201 {
                    print("User registered successfully")
                    return true
                } else {
                    let responseBody = String(data: data, encoding: .utf8) ?? "Unreadable response body"
                    print("Unexpected status code. Response body: \(responseBody)")
                    self.error = "Server responded with status code: \(httpResponse.statusCode)"
                    return false
                }
            } else {
                print("Invalid HTTPURLResponse")
                self.error = "Invalid response from server"
                return false
            }
        } catch {
            print("Exception occurred during request: \(error)")
            self.error = "Exception: \(error.localizedDescription)"
            return false
        }
    }

    func Login() async -> Bool {
        struct LoginResponse: Codable {
            let message: String
            let refreshToken: String
            let token: String
            let user: User
        }

        guard let url = URL(string: "\(Global.baseURL)/login") else {
            print("Invalid login URL")
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "email": user.email,
            "password": user.password ?? ""
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else {
            print("Failed to serialize JSON body")
            return false
        }

        request.httpBody = jsonData
        print("Sending request to \(url.absoluteString)")
        print("Request body: \(String(data: jsonData, encoding: .utf8) ?? "N/A")")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                print("Invalid response object")
                return false
            }

            print("Received response with status code: \(httpResponse.statusCode)")

            guard httpResponse.statusCode == 200 else {
                self.error = "Unexpected status code: \(httpResponse.statusCode)"
                print("Error: \(self.error ?? "Unknown error")")
                return false
            }

            print("Raw response data: \(String(data: data, encoding: .utf8) ?? "Unreadable")")

            let decoder = JSONDecoder()
            let loginResponse = try decoder.decode(LoginResponse.self, from: data)

            print("Decoded login response:")
            print("   • Message: \(loginResponse.message)")
            print("   • Token: \(loginResponse.token.prefix(20))...")
            print("   • Refresh Token: \(loginResponse.refreshToken.prefix(20))...")
            print("   • User ID: \(loginResponse.user.id)")

            UserDefaults.standard.setValue(loginResponse.token, forKey: "authToken")
            UserDefaults.standard.setValue(loginResponse.refreshToken, forKey: "refreshToken")

            UserDefaults.standard.removeObject(forKey: "currentUser")

            let encoder = JSONEncoder()
            let encodedUser = try encoder.encode(loginResponse.user)
            UserDefaults.standard.setValue(encodedUser, forKey: "currentUser")

            print("Saved user and tokens to UserDefaults")

            return true
        } catch {
            self.error = "Exception: \(error.localizedDescription)"
            print("Exception occurred: \(error)")
            return false
        }
    }
    
    func getAllUsers() async -> Bool {

        guard let url = URL(string: "\(Global.baseURL)/api/users") else {
            print("Invalid URL: \(Global.baseURL)/api/users")
            return false
        }

        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            print("Missing auth token from UserDefaults")
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        print("Sending GET request to \(url.absoluteString)")
        print("Authorization: Bearer \(token.prefix(20))...")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse {
                print("Received response with status code: \(httpResponse.statusCode)")

                if httpResponse.statusCode == 200 {
                    let decoder = JSONDecoder()
                    print("Raw response JSON: \(String(data: data, encoding: .utf8) ?? "Unreadable")")
                    let decodedUsers = try decoder.decode([User].self, from: data)
                    self.users = decodedUsers
                    print("Decoded \(decodedUsers.count) users")
                    return true
                } else {
                    let body = String(data: data, encoding: .utf8) ?? "Unreadable body"
                    print("Unexpected status code: \(httpResponse.statusCode)")
                    print("Response body: \(body)")
                    self.error = "Error: HTTP \(httpResponse.statusCode)"
                    return false
                }
            } else {
                print("Invalid HTTPURLResponse object")
                self.error = "Invalid response"
                return false
            }
        } catch {
            print("Request failed with error: \(error)")
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }
    
    func getUserByID() async -> Bool {

        guard let url = URL(string: "\(Global.baseURL)/api/users/\(user.id)") else {
            print("Invalid URL: \(Global.baseURL)/api/users/\(user.id)")
            return false
        }

        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            print("Missing auth token from UserDefaults")
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        print("Sending GET request to \(url.absoluteString)")
        print("Authorization: Bearer \(token.prefix(20))...")

        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse {
                print("Received response with status code: \(httpResponse.statusCode)")

                if httpResponse.statusCode == 200 {
                    let decoder = JSONDecoder()
                    print("Raw response JSON: \(String(data: data, encoding: .utf8) ?? "Unreadable")")
                    let decodedUser = try decoder.decode(User.self, from: data)
                    self.user = decodedUser
                    print("Decoded: \(decodedUser)")
                    return true
                } else {
                    let body = String(data: data, encoding: .utf8) ?? "Unreadable body"
                    print("Unexpected status code: \(httpResponse.statusCode)")
                    print("Response body: \(body)")
                    self.error = "Error: HTTP \(httpResponse.statusCode)"
                    return false
                }
            } else {
                print("Invalid HTTPURLResponse object")
                self.error = "Invalid response"
                return false
            }
        } catch {
            print("Request failed with error: \(error)")
            self.error = "Error: \(error.localizedDescription)"
            return false
        }
    }
    
    func updateUser(updateUser: User) async -> Bool {
        
        guard let url = URL(string: "\(Global.baseURL)/api/users") else {
            print("Invalid login URL")
            return false
        }
        
        guard let token = UserDefaults.standard.string(forKey: "authToken") else {
            self.error = "Missing auth token"
            print("Missing auth token from UserDefaults")
            return false
        }

        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        let body: [String: Any] = [
            "id": updateUser.email,
            "name": updateUser.name,
            "email": updateUser.email,
            "password": updateUser.password ?? "",
            "online": updateUser.online,
            "channels": updateUser.channels,
            "created": updateUser.created
        ]

        guard let jsonData = try? JSONSerialization.data(withJSONObject: body) else {
            print("Failed to serialize JSON body")
            return false
        }

        request.httpBody = jsonData
        print("Sending PUT request to \(url.absoluteString)")
        print("Authorization: Bearer \(token.prefix(20))...")
        
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            if let httpResponse = response as? HTTPURLResponse {
                print("Received response with status code: \(httpResponse.statusCode)")

                if httpResponse.statusCode == 200 {
                    print("User updated successfully")
                    return true
                } else {
                    let responseBody = String(data: data, encoding: .utf8) ?? "Unreadable response body"
                    print("Unexpected status code. Response body: \(responseBody)")
                    self.error = "Server responded with status code: \(httpResponse.statusCode)"
                    return false
                }
            } else {
                print("Invalid HTTPURLResponse")
                self.error = "Invalid response from server"
                return false
            }
        } catch {
            print("Exception occurred during request: \(error)")
            self.error = "Exception: \(error.localizedDescription)"
            return false
        }
        
    }

}
