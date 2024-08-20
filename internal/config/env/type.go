package env

const (
	runAddressKey       = "RUN_ADDRESS"
	dbURIKey            = "DATABASE_URI"
	accrualSystemURLKey = "ACCRUAL_SYSTEM_ADDRESS"
	jwtSecret           = "JWT_SECRET"
)

type Options struct {
	runAddress       string `env:"RUN_ADDRESS"`
	dbURI            string `env:"DATABASE_URI"`
	accrualSystemURL string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	jwtSecret        string `env:"JWT_SECRET"`
}

func (o *Options) HasRunAddress() bool {
	return o.runAddress != ""
}

func (o *Options) GetRunAddress() string {
	return o.runAddress
}

func (o *Options) HasDBUri() bool {
	return o.dbURI != ""
}

func (o *Options) GetDBUri() string {
	return o.dbURI
}

func (o *Options) HasAccrualSystemURL() bool {
	return o.accrualSystemURL != ""
}

func (o *Options) GetAccrualSystemURL() string {
	return o.accrualSystemURL
}

func (o *Options) HasJWTSecret() bool {
	return o.jwtSecret != ""
}

func (o *Options) GetJWTSecret() string {
	return o.jwtSecret
}
