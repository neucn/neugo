# NEU API

![Codecov](https://img.shields.io/codecov/c/github/neucn/neugo?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/neucn/neugo?style=flat-square)](https://goreportcard.com/report/github.com/neucn/neugo)
![Latest Tag](https://img.shields.io/github/v/tag/neucn/neugo?label=version&style=flat-square)


## 📈 Roadmap

- [x] Login
- [x] Extract Token
- [x] Encrypt WebVPN URL

## 🎨 Usage

```go
session := neugo.NewSession()
neugo.Use(session).WithAuth("student_id", "pass").Login(neugo.CAS)
neugo.Use(session).WithAuth("student_id", "pass").Login(neugo.WebVPN)
neugo.Use(session).WithToken("xxx").Login(neugo.CAS)
neugo.Use(session).WithToken("xxx").Login(neugo.WebVPN)

neugo.About(session).Token(neugo.CAS)
neugo.About(session).Token(neugo.WebVPN)

neugo.EncryptToWebVPN("http://ipgw.neu.edu.cn")
```

## 📃 License

MIT License.