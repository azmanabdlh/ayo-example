package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	logvar *log.Logger
)

func Init() {
	logvar = log.New(os.Stdout, "[LOGGER]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...any) {
	if logvar == nil {
		Init()
	}

	logvar.Output(2, "[INFO] "+fmt.Sprintf(format, v))
}
