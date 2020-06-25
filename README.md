# NEU API

🚧 WIP

## Usage

```go
session := neugo.NewSession()
neugo.Use(session).WithAuth("student_id","pass").On(neu.CAS).Login()
neugo.Use(session).WithAuth("student_id","pass").On(neu.WebVPN).Login()
neugo.Use(session).WithAuth("student_id","pass").On(neu.CAS).LoginService("xxx")
neugo.Use(session).WithAuth("student_id","pass").On(neu.WebVPN).LoginService("xxx")
neugo.Use(session).WithToken("xxx").On(neu.CAS).Login()
neugo.Use(session).WithToken("xxx").On(neu.WebVPN).Login()
neugo.Use(session).WithToken("xxx").On(neu.CAS).LoginService("xxx")
neugo.Use(session).WithToken("xxx").On(neu.WebVPN).LoginService("xxx")
```

## Roadmap

- [x] Login
- [x] Query Token
- [ ] Query Personal Info

## License

MIT License.