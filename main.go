package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type spigitCatalog struct {
	users []spigitUser
}

type spigitUser struct {
	username string
	email    string
}

type attCatalog struct {
	users []attUser
}

type attUser struct {
	attuid string
	cenet  string
}

func buildAttCatalog() (c attCatalog, err error) {
	csvfile, err := os.Open("/Users/jbarnett/Desktop/attextractpeople_simple.csv")

	if err != nil {
		return c, err
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	records, err := reader.ReadAll()
	if err != nil {
		return c, err
	}

	for _, row := range records {
		cenet := row[0]
		attuid := row[1]
		attuser := attUser{attuid: attuid, cenet: cenet}
		c.users = append(c.users, attuser)
	}

	return c, nil
}

func buildspigitCatalog() (c spigitCatalog, err error) {
	csvfile, err := os.Open("/Users/jbarnett/Desktop/tip1_users_simple.csv")

	if err != nil {
		return c, err
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	//reader.FieldsPerRecord = 0
	records, err := reader.ReadAll()
	if err != nil {
		return c, err
	}

	for _, row := range records {
		username := row[0]
		email := row[1]
		user := spigitUser{username: username, email: email}
		c.users = append(c.users, user)
	}

	return c, nil
}

func (c *attCatalog) findByAttuid(attuid *string) (attUser, error) {
	for _, user := range c.users {
		if user.attuid == *attuid {
			return user, nil
		}
	}
	return attUser{}, errors.New("Cannot find user!")
}

func findEmptyEmails(sc *spigitCatalog, ac *attCatalog) {
	for i, user := range sc.users {
		attuser, err := ac.findByAttuid(&user.username)
		if err != nil {
			// fmt.Printf("[%v] orphaned spigit engage user: %+v\n", i, user)
		} else {
			if attuser.cenet == "" {
				fmt.Printf("[%v] %+v\n", i, attuser)
			}
		}
	}
}

func main() {
	ac, err := buildAttCatalog()
	if err != nil {
		fmt.Println(err)
		return
	}

	sc, err := buildspigitCatalog()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("ATT Catalog Users: %v\n", len(ac.users))
	fmt.Printf("Spigit Engage Users: %v\n", len(sc.users))
	findEmptyEmails(&sc, &ac)
}
