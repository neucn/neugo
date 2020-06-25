# NEU API

🚧 WIP

## Usage

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
```

## Roadmap

- [x] Login
- [x] Query Token
- [ ] Query Personal Info

## License

MIT License.