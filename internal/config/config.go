package config

import (
	"html/template"
	"log"

	"github.com/HINKOKO/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}

// beware of IMPORT CYCLE !!!
