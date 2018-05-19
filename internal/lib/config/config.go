package config

// CertList maps the certificate list in config to a struct
type CertList struct {
	SubjectCN        string   `mapstructure:"cn"`
	KeyAlgo          string   `mapstructure:"keyalgo"`
	PrivKey          string   `mapstructure:"privkey"`
	Cert             string   `mapstructure:"cert"`
	Chain            string   `mapstructure:"chain"`
	FullChain        string   `mapstructure:"fullchain"`
	FullChainPrivKey string   `mapstructure:"fullchainprivkey"`
	Hooks            []string `mapstructure:"hooks"`
}
