package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

//Email representa a estrutura de um email
type Email struct {
	Email        string
	_password    string
	_cryptMethod string
}

func (e *Email) makePassword() string {
	return fmt.Sprintf("%v%v", e._cryptMethod, e._password)
}

var (
	domain     = flag.String("d", "", "domain name, ex: domain.com")
	shadowPath = flag.String("f", "", "Path to shadow file")
	emails     *[]Email
)

func init() {
	flag.Parse()
}

func main() {
	if len(*shadowPath) == 0 || len(*domain) == 0 {
		log.Fatal("Need arguments use --help")
	}
	prepareEmails()
	createSQL()
}

func createSQL() {
	var data string
	for _, ac := range *emails {
		fmt.Printf("Making mysql script of %+v \n", ac.Email)
		data += fmt.Sprintf("UPDATE mailbox SET password = '%v' WHERE username = '%v';\n", ac.makePassword(), ac.Email)
	}
	saveToFile(data)
}

func saveToFile(data string) {
	fmt.Print("saving to file...")
	f, err := os.Create("emails.sql")
	if err != nil {
		panic(err)
	}
	f.Write([]byte(data))
	f.Sync()
}

func prepareEmails() {
	rawEmails := readFile()
	if len(*rawEmails) == 0 {
		log.Fatal("File empty!")
	}
	checkDataStruct(rawEmails)
	buildEmails(rawEmails)
}

func buildEmails(rawEmails *[]string) {
	var x []Email
	for _, rawEmail := range *rawEmails {
		splitedEmail := strings.Split(rawEmail, ":")
		email := Email{
			Email:        fmt.Sprintf("%v@%v", splitedEmail[0], *domain),
			_password:    splitedEmail[1],
			_cryptMethod: defineCrypt(rawEmail),
		}
		x = append(x, email)
	}
	emails = &x
}

func checkDataStruct(emails *[]string) {
	for _, email := range *emails {
		qtd := strings.Split(email, ":")
		if len(qtd) != 9 {
			log.Fatal("Data struct error")
		}
	}
}

func defineCrypt(email string) string {
	d := strings.Split(email, ":")[1]
	switch d[:3] {
	case "$1$":
		return "{MD5-CRYPT}"
	case "$5$":
		return "{SHA256-CRYPT}"
	case "$6$":
		return "{SHA512-CRYPT}"
	default:
		return "{SHA512-CRYPT}"
	}
}

func readFile() (emails *[]string) {
	var e []string
	file, err := os.Open(*shadowPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		e = append(e, scanner.Text())
	}
	emails = &e

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
