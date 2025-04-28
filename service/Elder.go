package service

type Elder struct{}

func (Elder) Register() {}

func (Elder) Update() {}

func (Elder) Get() {}

func (Elder) Join() {}

func (Elder) Leave() {}

// 决定是否通过
func (Elder) Decide() {}
