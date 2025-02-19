package core

import (
	"log/slog"
	"net/smtp"

	"cloud.google.com/go/vertexai/genai"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

type CoreOpts struct {
	Port       int
	MailAddr   string
	MailSecret string
	SmtpAddr   string
	MongoDbUri string

	ModelName    string
	GcpProjectID string
	GcpLocation  string
}

// Core defines the core construct of the service.
type Core struct {
	Port int
	opts *CoreOpts

	smtpAuth smtp.Auth
	llm      *genai.GenerativeModel
	DB       *repository.MongoDBClient
	Lo       *slog.Logger
}

// configureIndexesDB helps to configure database level constraints and checks.
// If complaints and rules on the collections as per the defined actions to make data integrity stronger.
func (co *Core) configureIndexesDB() {
	co.configureUsersIndexes()
}

func NewCore(opts *CoreOpts) *Core {

	co := &Core{
		opts: opts,
		Port: opts.Port,
		Lo:   slog.Default(),
	}

	// Initialize the database
	mdb, err := intiializeDatabase(co.opts.MongoDbUri)
	if err != nil {
		co.Lo.Error("not able to connect to mongoDB", slog.Any("mdb_err", err.Error()))
		panic(err)
	}
	co.DB = &repository.MongoDBClient{
		Client: mdb,
	}

	// Initialize the SMTP Auth instance to be reused.
	co.smtpAuth = smtp.PlainAuth(
		"",
		co.opts.MailAddr,
		co.opts.MailSecret,
		co.opts.SmtpAddr,
	)

	// Configure indexes managements
	go co.configureIndexesDB()

	// Initialize LLM model (Gemini).
	if err := co.initializeLLM(); err != nil {
		panic(err)
	}

	return co
}
