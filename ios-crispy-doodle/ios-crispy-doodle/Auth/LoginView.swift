//
//  LoginView.swift
//  ios-crispy-doodle
//
//  Created by Peter Bishop on 5/12/25.
//

import SwiftUI

struct LoginView: View {
    
    @StateObject var userVM = UserViewModel()
     @State var password: String = ""
     @State var newUser: Bool = false
     @State var showAlert: Bool = false
     @State var loginSuccess: Bool = false
    
    var body: some View {
        NavigationStack{
                           VStack{
                               Spacer()
                               Text("Login").font(.system(size: 34))
                                   .fontWeight(.ultraLight)
                               Divider().padding()
                               TextField("Email", text: $userVM.user.email)
                                   .tint(.black)
                                   .autocapitalization(.none)
                                   .disableAutocorrection(true)
                                   .padding()
                               
                               SecureField("Password", text:  $password)
                                   .tint(.black)
                                   .autocapitalization(.none)
                                   .disableAutocorrection(true)
                                   .padding()
                               
                               
                               Button("Submit", action: {
                                   userVM.user.password = password
                                   Task{
                                       let result = await userVM.Login()
                                       await MainActor.run {
                                           loginSuccess = result
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
                               .onChange(of: loginSuccess, {
                                   oldResponse, newResponse in
                                   if !newResponse {
                                       showAlert = true
                                   }
                               })
                               .alert("Error", isPresented: $showAlert) {
                                               Button("OK", role: .cancel) {
                                                   userVM.user.email = ""
                                                   userVM.user.password = ""
                                               }
                                           } message: {
                                               Text(String(userVM.error ?? "Error"))
                                           }
                                           .navigationDestination(isPresented: $loginSuccess, destination: {
                                               SuccessView().navigationBarBackButtonHidden(true)
                                           })
                               Spacer()
                               HStack{
                                   Spacer()
                                   Text("I don't have an account.").fontWeight(.ultraLight)
                                   Button("Register", action: {
                                       newUser = true
                                   }).foregroundStyle(.black)
                                       .fontWeight(.light)
                                       .navigationDestination(isPresented: $newUser, destination: {
                                           RegisterView().navigationBarBackButtonHidden(true)
                                       })
                                   Spacer()
                               }
                           }.onAppear{
                               
                           }
                       }
    }
}

#Preview {
    LoginView()
}
