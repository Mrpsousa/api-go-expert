package log

import (
	"log"
	"os"
	"time"
)

// TODO:_ criar log de cor vermelha pra erro, amarelo pra warr e branco  pra log normais

func LogMessage(name, text string, saveToFile bool) {
	// Configuração do log para incluir data e hora
	log.SetFlags(log.Ldate | log.Ltime)

	// Se saveToFile for verdadeiro, abre um arquivo para registrar os logs
	var logFile *os.File
	if saveToFile {
		var err error
		logFile, err = os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(logFile)
			defer logFile.Close()
		}
	}

	// Obtém a data e hora atual
	currentTime := time.Now()

	// Formata a mensagem
	message := currentTime.Format("2006-01-02 15:04:05") + " - " + name + ": " + text

	// Registra a mensagem no log
	log.Println(message)
}

// logMessage("Usuário1", "Este é um exemplo de log.", true)
