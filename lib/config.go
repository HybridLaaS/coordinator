package lib

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Environment struct {
	// Server setup
	Host   string `env:"HOST,default=0.0.0.0"`
	Port   int    `env:"PORT,default=8090"`
	TlsDir string `env:"TLS_DIR"`

	// Database setup
	DBFile      string `env:"DB_FILE,default=opnlaas.db"`
	DBSalt      string `env:"DB_SALT,required=true"`
	DBQueueSize int    `env:"DB_QUEUE_SIZE,default=256"`

	// SMTP Email setup
	SmtpHost     string `env:"SMTP_HOST,required=true"`
	SmtpPort     int    `env:"SMTP_PORT,required=true"`
	SmtpUser     string `env:"SMTP_USER,required=true"`
	SmtpPassword string `env:"SMTP_PASSWORD,required=true"`

	// Configuration
	LabName              string   `env:"LAB_NAME,default=Sample Laboratory"`
	LabOrg               string   `env:"LAB_ORG,default=Placebo Pharmaceuticals"`
	LabContact           string   `env:"LAB_CONTACT,required=true"`
	EmailDomainWhiteList []string `env:"EMAIL_DOMAIN_WHITELIST,default=example.com|placebo.org"`
}

var Config Environment

func InitEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	_, err := env.UnmarshalFromEnviron(&Config)
	if err != nil {
		return err
	}

	return nil
}
