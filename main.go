package main

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"log"
	"os"
)

func main() {
	host := os.Getenv("SFTP_HOST")
	port := os.Getenv("SFTP_PORT")
	user := os.Getenv("SFTP_USER")
	addr := fmt.Sprintf("%s:%s", host, port)

	privateKeyB64 := os.Getenv("SFTP_PRIVATE_KEY")
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyPass := os.Getenv("SFTP_PRIVATE_KEY_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(privateKeyPass))
	if err != nil {
		fmt.Println("alpha")
		log.Fatal(err)
	}

	auth := []ssh.AuthMethod{
		ssh.PublicKeys(signer),
	}

	knownHosts := os.Getenv("KNOWN_HOSTS_PATH")
	hostKey, err := knownhosts.New(knownHosts)
	if err != nil {
		fmt.Println("beta")
		log.Fatal(err)
	}

	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: hostKey,
	}

	conn, err := ssh.Dial("tcp", addr, conf)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
}
