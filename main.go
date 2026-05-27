package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Site struct {
	URL string
	Intervalo int
}

func main() {
	for {
		menu()

		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			resposta := strings.TrimSpace(scanner.Text())
			respostaConvertida, err := strconv.Atoi(resposta)

			if err != nil {
				fmt.Println("Erro ", err.Error())
			}

			if respostaConvertida == 1 {
				limpaTerminal()

    			sites := listaSites()

				// Cria um channel que transporta as strings dos sites
				ch := make(chan Site)

				for _, site := range sites {
					go verificaSite(site, ch)
				}

				for {
					// Espera receber um valor do channel criado
					site := <-ch
				
					// Espera 60 segundos antes de verificar novamente
					go func(site Site) {
						time.Sleep(time.Duration(site.Intervalo) * time.Second)
					
						// Verifica novamente criando nova goroutine
						verificaSite(site, ch)
					} (site)
				}
			} else if respostaConvertida == 2 {
				limpaTerminal()
				sites := listaSites()
				for _, site := range sites {
					fmt.Println(site)
				}
			} else if respostaConvertida == 3 {
				limpaTerminal()
				fmt.Println("\nPara cadastrar novo site, digite [site],[intervalo]")
				scannerSite := bufio.NewScanner(os.Stdin)

				if scannerSite.Scan() {
					resposta := strings.Split(strings.TrimSpace(scannerSite.Text()), ",")

					if len(resposta) != 2 {
						fmt.Println("Formato inválido.")
						return
					}

					respostaArr := [2]string{resposta[0], resposta[1]}

					cadastraSite("sitesMonitorados.txt", respostaArr)
				}
			} else {
				os.Exit(0)
			}
		}
	}
}

func verificaSite(url Site, ch chan Site) {
	defer func() {
		ch <- url
	} ()
	// Faz requisição HTTP no site com timeout:
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url.URL)

	if err != nil {
		registraErros(url, 0)

		// Enviar o site de volta para o channel
		return
	}

	// Fecha a resposta para liberar os recursos
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("Site online: ", url.URL)
	} else {
		registraErros(url, resp.StatusCode)
	}
}

func registraErros(url Site, statusCode int) {
	arquivo, err := os.OpenFile("logserr.txt", os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

	check(err)

	defer func() {
		if err := arquivo.Close(); err != nil {
			check(err)
		}
	}()

	timeNow := time.Now().Format("2006-01-02 15:04:05")

	log := fmt.Sprintf("[%s]: %s | HTTP: %d\n", timeNow, url.URL, statusCode)
	_, err = arquivo.WriteString(log)
	check(err)
}

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func listaSites() []Site {
	arquivo, err := os.Open("sitesMonitorados.txt")
	check(err)

	defer arquivo.Close()
	
	var sites []Site
	leitor := bufio.NewScanner(arquivo)

	if leitor.Err() != nil {
		fmt.Println("Erro no leitor", leitor.Err().Error())
	}

	for leitor.Scan() {
		linha := leitor.Text()

		partes := strings.Split(linha, ",")
		url := partes[0]
		intervalo, err := strconv.Atoi(partes[1])
		check(err)

		site := Site{
			URL: url,
			Intervalo: intervalo,
		}

		sites = append(sites, site)
	}

	return sites
}

func cadastraSite(nomeArq string, conteudo [2]string) {
	file, err := os.OpenFile(nomeArq, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Erro ao cadastrar site", err.Error())
		return
	}

	defer file.Close()

	dadoSite := conteudo[0]
	dadoIntervalo := conteudo[1]

	dadosConcatenados := dadoSite + "," + dadoIntervalo + "\n"

	file.Write([]byte(dadosConcatenados))
}

func menu() {
	fmt.Printf("Monitoramento - Sookup (Simple Lookup)\n")
	fmt.Println("1 - Iniciar monitoramento\n2 - Listar sites monitorados\n3 - Cadastrar novo site\n0 - Sair")
}

func limpaTerminal() {
	cmd := exec.Command("cmd", "clear")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	check(err)
}