package config

const (
	// Dev 開発環境
	Dev Env = "dev"
	// Pro 本番環境
	Pro Env = "pro"
	// Test テスト環境
	Test Env = "test"
)

// Env 実行環境 (dev, pro, test)
type Env string

// IsDev 開発環境か返す
func (e *Env) IsDev() bool {
	if !e.IsInitialized() {
		panic("Env is not initialized")
	}
	return (*e) == Dev
}

// IsPro 本番環境か返す
func (e *Env) IsPro() bool {
	if !e.IsInitialized() {
		panic("Env is not initialized")
	}
	return (*e) == Pro
}

// IsTest テスト環境か返す
func (e *Env) IsTest() bool {
	if !e.IsInitialized() {
		panic("Env is not initialized")
	}
	return (*e) == Test
}

// ToString stringで返す
func (e *Env) ToString() string {
	return string(*e)
}

// IsInitialized 初期化済みか返す
func (e *Env) IsInitialized() bool {
	switch *e {
	case Dev, Pro, Test:
		return true
	}
	return false
}
