package model

type Config struct {
	CommonName         string   `yaml:"commonName"`
	OrganizationalUnit []string `yaml:"organizationalUnit"`
	Organization       []string `yaml:"organization"`
	Country            []string `yaml:"country"`
	DNSNames           []string `yaml:"dnsNames"`
	Exp                int      `yaml:"exp"`
}
