// Package neugo
//
// This package encapsulates some operations on the campus SSO service, in particular CAS, of NEU (cn).
//
// Examples:
//
// E1. Log in to the CAS using account and get the token.
//   // you can use your own *http.Client instead of creating a new session.
//   client := neugo.NewSession()
//   err := neugo.Use(client).WithAuth("student_id", "password").Login(CAS)
//   token := neugo.About(client).Token(CAS)
//
// E2. Request a service via Web VPN using token.
//   client := neugo.NewSession()
//   err := neugo.Use(client).WithToken("your_webvpn_token").Login(WebVPN)
//   serviceURL := "https://ipgw.neu.edu.cn"
//	 reqURL := neugo.EncryptURLToWebVPN(serviceURL)
//   req := http.NewRequest("GET", reqURL, nil)
//   resp, err := client.Do(req)
//
package neugo
