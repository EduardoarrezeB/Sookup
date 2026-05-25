package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"strconv"
)

type Site struct {
	URL string
	Intervalo int
}

func main() {
	// Lista de sites para monitorar
    sites := listaSites()

	// Cria um channel que transporta as strings dos sites
	ch := make(chan Site)

	for _, site := range sites {
		go verificaSite(site, ch)
	}

	// Continua monitorando
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
}

func verificaSite(url Site, ch chan Site) {
	defer func() {
		ch <- url
	} ()
	// Faz requisição HTTP no site com timeout:
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url.URL)

	if err != nil {
		check(err)

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

	horario := time.Now().Format("2006-01-02 15:04:05")

	log := fmt.Sprintf("%s: Problema no site: %s | HTTP: %d\n", horario, url.URL, statusCode)
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