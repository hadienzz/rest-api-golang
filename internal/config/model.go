package config

type Config struct {
	Port               string
	JwtKey             string
	DatabaseDirectURL  string
	DatabasePoolerURL  string
	SupabaseURL        string
	SupabaseServiceKey string
	MidtransClientKey  string
	MidtransServerKey  string
	MidtransMerchantID string
	DatabaseURL        string
}
