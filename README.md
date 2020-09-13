# NEU API

![Codecov](https://img.shields.io/codecov/c/github/neucn/neugo?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/neucn/neugo?style=flat-square)](https://goreportcard.com/report/github.com/neucn/neugo)
![Latest Tag](https://img.shields.io/github/v/tag/neucn/neugo?label=version&style=flat-square)

🚧 WIP

## 📈 Roadmap

- [x] Login
- [x] Query Token
- [ ] Query Personal Info

## 🎨 Usage

```go
session := neugo.NewSession()
neugo.Use(session).WithAuth("student_id","pass").On(neugo.CAS).Login()
neugo.Use(session).WithAuth("student_id","pass").On(neugo.WebVPN).Login()
neugo.Use(session).WithAuth("student_id","pass").On(neugo.CAS).LoginService("xxx")
neugo.Use(session).WithAuth("student_id","pass").On(neugo.WebVPN).LoginService("xxx")
neugo.Use(session).WithToken("xxx").On(neugo.CAS).Login()
neugo.Use(session).WithToken("xxx").On(neugo.WebVPN).Login()
neugo.Use(session).WithToken("xxx").On(neugo.CAS).LoginService("xxx")
neugo.Use(session).WithToken("xxx").On(neugo.WebVPN).LoginService("xxx")

neugo.About(session).Token(neugo.WebVPN)
neugo.About(session).Token(neugo.CAS)
```

## 📃 License

MIT License.