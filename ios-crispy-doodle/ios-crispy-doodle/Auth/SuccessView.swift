//
//  SuccessView.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import SwiftUI

// gets user stored in UserDefaults
// gets tokens stored in UserDefaults
// gets all users

struct SuccessView: View {
    @StateObject private var userVM = UserViewModel()

    @State private var currentUser: User?
    @State private var jwt: String?
    @State private var refreshToken: String?
    @State private var errorLoading: Bool = false
    @State private var showProfile: Bool = false
    @State private var foundUsers: Bool = false
    @State private var logout: Bool = false

    private func logoutUser() {
        UserDefaults.standard.removeObject(forKey: "currentUser")
        UserDefaults.standard.removeObject(forKey: "authToken")
        UserDefaults.standard.removeObject(forKey: "refreshToken")
        logout = true
    }

    private var logoutButton: some View {
        Button("Logout", action: logoutUser)
            .fontWeight(.ultraLight)
            .foregroundColor(.black)
            .padding()
            .background(
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.white)
                    .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
            )
    }
    
    private var refreshButton: some View {
        Button("Refresh", action: {
            Global.refreshAccessToken { result in
                switch result {
                case .success(let newToken):
                    print("Refreshed token: \(newToken)")
                case .failure(let error):
                    print("Failed to refresh: \(error.localizedDescription)")
                }
            }
        })
            .fontWeight(.ultraLight)
            .foregroundColor(.black)
            .padding()
            .background(
                RoundedRectangle(cornerRadius: 8)
                    .fill(Color.white)
                    .shadow(color: .gray.opacity(0.4), radius: 4, x: 2, y: 2)
            )
    }

    var body: some View {
        NavigationStack {
            VStack {
                if errorLoading {
                    Text("Failed to load user info.")
                        .foregroundColor(.red)
                        .padding()
                }

                if userVM.isLoading {
                    ProgressView("Loading...")
                        .padding()
                    logoutButton
                } else {
                    HStack {
                        Text(currentUser?.name ?? "No user")
                        Spacer()
                        logoutButton
                    }
                    .padding()
                    NavigationView {
                        Group {
                            if userVM.isLoading {
                                ProgressView("Loading Users...")
                            } else if userVM.error != nil {
                                Text("Error: \(String(describing: userVM.error))")
                                    .foregroundColor(.red)
                            } else {
                                ScrollView {
                                    LazyVStack(spacing: 12) {
                                        ForEach(userVM.users) { user in
                                            if user.id != currentUser?.id {
                                            HStack{
                                                    Button {
                                                        showProfile = true
                                                    } label: {
                                                        Image(systemName: "info.circle.fill")
                                                            .tint(.black)
                                                    }
                                                    .sheet(isPresented: $showProfile) {
                                                        VStack(alignment: .leading, spacing: 16) {
                                                            HStack {
                                                                Text(user.name).font(.title2).bold()
                                                                if user.online {
                                                                    Image(systemName: "checkmark.circle.fill")
                                                                        .foregroundColor(.green)
                                                                }
                                                            }
                                                            Text(user.email).font(.body)
                                                            Spacer()
                                                        }
                                                        .padding()
                                                    }
                                                    VStack(alignment: .leading, spacing: 2) {
                                                        Text(user.name)
                                                            .font(.headline)
                                                        Text(user.email)
                                                            .font(.subheadline)
                                                            .foregroundColor(.secondary)
                                                    }
                                                    Spacer()
                                                    Button {
                                                        // Create Channel
                                                        // Update Users with Channel ID
                                                    } label: {
                                                        Image(systemName: "chevron.right").tint(.black)
                                                    }
                                                }
                                                .frame(maxWidth: .infinity, alignment: .leading)
                                                .padding()
                                                .background(Color.white)
                                                .cornerRadius(12)
                                                .shadow(color: Color.black.opacity(0.1), radius: 4, x: 0, y: 2)
                                            }
                                        }
                                    }
                                }.padding()

                            }
                        }
                        .navigationTitle("Users")
                    }
                }
            }
            .navigationDestination(isPresented: $logout) {
                LoginView().navigationBarBackButtonHidden(true)
            }
            .onAppear {
                do {
                    
                    if let data = UserDefaults.standard.data(forKey: "currentUser") {
                        do {
                            print(String(data: data, encoding: .utf8) ?? "No data")
                            currentUser = try User.decodeUser(from: data)
                            currentUser?.online = true
                            Task {
                                try await userVM.updateUser(updateUser: currentUser!)
                            }
                        } catch {
                            print("Failed to decode user: \(error.localizedDescription)")
                            currentUser = nil // or fallback User()
                        }
                    } else {
                        print("No user data found in UserDefaults")
                        currentUser = nil // or fallback User()
                    }
                    jwt = UserDefaults.standard.string(forKey: "authToken")
                    refreshToken = UserDefaults.standard.string(forKey: "refreshToken")
                    Task {
                        await userVM.getAllUsers()
                    }
                }
            }
        }
    }
}

#Preview {
    SuccessView()
}
