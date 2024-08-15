package env

const (
	runAddressKey       = "RUN_ADDRESS"
	dbUriKey            = "DATABASE_URI"
	accrualSystemUrlKey = "ACCRUAL_SYSTEM_ADDRESS"
)

type Options struct {
	runAddress       string `env:"RUN_ADDRESS"`
	dbUri            string `env:"DATABASE_URI"`
	accrualSystemUrl string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (o *Options) HasRunAddress() bool {
	return o.runAddress != ""
}

func (o *Options) GetRunAddress() string {
	return o.runAddress
}

func (o *Options) HasDBUri() bool {
	return o.dbUri != ""
}

func (o *Options) GetDBUri() string {
	return o.dbUri
}

func (o *Options) HasAccrualSystemUrl() bool {
	return o.accrualSystemUrl != ""
}

func (o *Options) GetAccrualSystemUrl() string {
	return o.accrualSystemUrl
}
