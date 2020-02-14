package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"time"

	// "github.com/catalinaM/go-whatsapp/binary/proto"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/catalinaM/go-whatsapp"
)

func main() {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	// <-time.After(3 * time.Second)
	ch, err := wac.GetGroupMetaData("@g.us")
 var response map[string]interface{}

	err = json.Unmarshal([]byte(<-ch), &response)

	if err != nil {
		fmt.Println("error decoding response message: %v\n", err)
	}


	fmt.Println(response)
	wac.GroupAnnouceFlag("@g.us", "@c.us", true)
// 	var response1 map[string]interface{}
// 	fmt.Println("back")
// 	 err = json.Unmarshal([]byte(<-ch1), &response1)
// fmt.Println("gata")
// 	 if err1 != nil {
// 		 fmt.Println("error decoding response message: %v\n", err)
// 	 }

	 fmt.Println(response1)
}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}