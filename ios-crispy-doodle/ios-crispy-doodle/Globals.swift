//
//  Globals.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import Foundation

struct Global {
    
    static var baseURL: String = "http://localhost:8080"
    
    struct TokenResponse: Codable {
        let access_token: String
    }

    static func refreshAccessToken(completion: @escaping (Result<String, Error>) async -> Void) {
        guard let refreshToken = UserDefaults.standard.string(forKey: "refreshToken") else {
            Task {
                await completion(.failure(NSError(domain: "", code: 401, userInfo: [NSLocalizedDescriptionKey: "Missing refresh token"])))
            }
            return
        }

        guard let url = URL(string: "\(baseURL)/api/refresh") else {
            Task {
                await completion(.failure(NSError(domain: "", code: 400, userInfo: [NSLocalizedDescriptionKey: "Invalid refresh URL"])))
            }
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let requestBody = ["refresh_token": refreshToken]
        do {
            request.httpBody = try JSONEncoder().encode(requestBody)
        } catch {
            Task {
                await completion(.failure(error))
            }
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                Task {
                    await completion(.failure(error))
                }
                return
            }

            guard let httpResponse = response as? HTTPURLResponse, httpResponse.statusCode == 200 else {
                Task {
                    await completion(.failure(NSError(domain: "", code: 401, userInfo: [NSLocalizedDescriptionKey: "Unauthorized"])))
                }
                return
            }

            guard let data = data else {
                Task {
                    await completion(.failure(NSError(domain: "", code: 500, userInfo: [NSLocalizedDescriptionKey: "No data received"])))
                }
                return
            }

            do {
                let decoded = try JSONDecoder().decode(TokenResponse.self, from: data)
                // Save the new token
                UserDefaults.standard.set(decoded.access_token, forKey: "authToken")
                Task {
                    await completion(.success(decoded.access_token))
                }
            } catch {
                Task {
                    await completion(.failure(error))
                }
            }
        }.resume()
    }

}

extension NSError: @retroactive LocalizedError {
    
    public var errorDescription: String? { localizedDescription }
    
    public var failureReason: String? { localizedFailureReason }
    
    public var recoverySuggestion: String? { localizedRecoverySuggestion }
    
    // public var helpAnchor: String? { get }  // Goes as is.
}
