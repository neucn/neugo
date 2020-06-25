package neugo

// 个人信息
type PersonalInfo struct {
	Profile profile
	Mail    mail
	Wallet  wallet
	Network network
	Library library
}

// 个人基本信息
type profile struct {
	StudentID string
	Name      string
	Gender    string
	College   string
	Role      string
}

// 学校邮箱
type mail struct {
	Unread int
	Total  int
}

// 校园卡钱包
type wallet struct {
	Balance float64
	Subsidy float64
}

// 校园网
type network struct {
	Balance     float64
	UsedTraffic string
}

// 图书外借
type library struct {
	Total   int
	Current int
}

func parseInfo(body string) (*PersonalInfo, error) {
	// TODO
	return nil, nil
}
