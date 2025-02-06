package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Verificacao struct {
	Timestamp time.Time `json:"timestamp"`
	Caminho   string    `json:"caminho"`
	Verificado bool     `json:"verificado"`
	Resposta  string    `json:"resposta"`
	Erro      string    `json:"erro,omitempty"`
	Host      string    `json:"host"`
}

func verificarApp(caminho string) Verificacao {
	cmd := exec.Command("codesign", "--verify", "--deep", "--strict", "--verbose=2", caminho)
	resposta, err := cmd.CombinedOutput()
	host, _ := os.Hostname() // Obtém o hostname da máquina
	resultado := Verificacao{
		Timestamp: time.Now(),
		Caminho:   caminho,
		Verificado: err == nil,
		Resposta:   string(resposta),
		Host:       host,
	}

	if err != nil {
		resultado.Erro = fmt.Sprintf("verificação falhou: %s", err)
		log.Printf("Verificação falhou para %s: %v\nResposta: %s\n", caminho, err, string(resposta))
	} else {
		log.Printf("Verificação bem-sucedida para %s: %s\n", caminho, string(resposta))
	}

	return resultado
}

func registrarResultadoComoJSON(resultado Verificacao) {
	dadosLog, err := json.Marshal(resultado)
	if err != nil {
		log.Fatalf("Falha ao converter resultado: %v", err)
	}

	// Exibe o log no formato JSON para ser capturado
	fmt.Println(string(dadosLog))
}

func main() {
	aplicativos := []string{
		"/Applications/Google Chrome.app",
		"/Applications/Safari.app",
	}

	for _, caminho := range aplicativos {
		resultado := verificarApp(caminho)
		registrarResultadoComoJSON(resultado)
	}
}
