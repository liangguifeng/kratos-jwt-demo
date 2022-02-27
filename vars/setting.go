package vars

// Inner Api相关配置
type JwtSettings struct {
	SecretKey string
}

var JwtSetting = &JwtSettings{}
