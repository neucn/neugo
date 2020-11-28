package neugo

import "fmt"

// 使用示例
func Example() {
	client := NewSession()
	// 登录
	err := Use(client).WithAuth("student_id", "password").Login(CAS)
	if err != nil {
		panic(err)
	}

	// 提取 Token
	token, err := About(client).Token(CAS)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	// 将服务 URL 转换为 WebVPN 上对应的 URL
	url := EncryptToWebVPN("http://ipgw.neu.edu.cn")
	fmt.Println(url)
}
