package service

type Organization struct{}

func (Organization) Register() {}

func (Organization) Update() {}

func (Organization) Get() {}

// 决定通过志愿者，管理员，老人
func (Organization) Decide() {}
