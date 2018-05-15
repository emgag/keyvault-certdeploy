package config

type CertList struct {
	SubjectCN string   `mapstructure:"cn"`
	KeyAlgo   string   `mapstructure:"keyalgo"`
	PrivKey   string   `mapstructure:"privkey"`
	Cert      string   `mapstructure:"cert"`
	FullChain string   `mapstructure:"fullchain"`
	Hooks     []string `mapstructure:"hooks"`
}
