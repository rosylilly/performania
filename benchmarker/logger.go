package benchmarker

import (
	"log"
	"os"
)

var (
	UserLogger = log.New(os.Stdout, "[USER] ", log.Ltime|log.Lmicroseconds)
	DevLogger  = log.New(os.Stderr, "[DEV ] ", log.Ltime|log.Lmicroseconds)
)
