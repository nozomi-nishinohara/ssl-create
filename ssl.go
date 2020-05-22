package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/nozomi-nishinohara/ssl-create/model"
	"gopkg.in/yaml.v2"
)

var (
	config model.Config
)

func init() {
	var (
		configName = flag.String("setting", "config.yaml", "Configファイル名")
	)
	flag.Parse()
	byteConfig, err := ioutil.ReadFile(*configName)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(byteConfig, &config); err != nil {
		panic(err.Error())
	}
}

func save(filename string, block *pem.Block) {
	f, _ := os.Create(filename)
	defer f.Close()
	pem.Encode(f, block)
	f.Close()
}

func ca() (caTpl *x509.Certificate, privateCaKey *rsa.PrivateKey) {
	privateCaKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	publicCaKey := privateCaKey.Public()

	subjectCa := pkix.Name{
		CommonName:         config.CommonName,
		OrganizationalUnit: config.OrganizationalUnit,
		Organization:       config.Organization,
		Country:            config.Country,
	}
	caTpl = &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               subjectCa,
		NotAfter:              time.Now().AddDate(10, 0, 0),
		NotBefore:             time.Now(),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}
	caCertificate, _ := x509.CreateCertificate(rand.Reader, caTpl, caTpl, publicCaKey, privateCaKey)
	block := &pem.Block{Type: "CERTIFICATE", Bytes: caCertificate}
	save("ssl/ca.crt", block)
	return
}

func serverSSL(caTpl *x509.Certificate, privateCaKey *rsa.PrivateKey) {
	privateSslKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicSslKey := privateSslKey.Public()
	subjectSsl := pkix.Name{
		CommonName:         config.CommonName,
		OrganizationalUnit: config.OrganizationalUnit,
		Organization:       config.Organization,
		Country:            config.Country,
	}
	sslTpl := &x509.Certificate{
		SerialNumber: big.NewInt(123),
		Subject:      subjectSsl,
		NotAfter:     time.Now().AddDate(10, 0, 0),
		NotBefore:    time.Now(),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     config.DNSNames,
	}
	derSslCertificate, _ := x509.CreateCertificate(rand.Reader, sslTpl, caTpl, publicSslKey, privateCaKey)
	block := &pem.Block{Type: "CERTIFICATE", Bytes: derSslCertificate}
	save("ssl/server.crt", block)
	derPrivateSslKey := x509.MarshalPKCS1PrivateKey(privateSslKey)
	block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derPrivateSslKey}
	save("ssl/server.key", block)
}
