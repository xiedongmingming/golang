package main

import (
	"log"
	"github.com/jtblin/go-ldap-client"
)

func main() {

	client := &ldap.LDAPClient{
		Base:         "dc=example,dc=com",
		Host:         "ldap.example.com",
		Port:         389,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=People,dc=example,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}

	// It is the responsibility of the caller to close the connection

	defer client.Close()

	ok, user, err := client.Authenticate("username", "password")

	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "username", err)
	}

	if !ok {
		log.Fatalf("Authenticating failed for user %s", "username")
	}

	log.Printf("User: %+v", user)

	groups, err := client.GetGroupsOfUser("username")

	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "username", err)
	}

	log.Printf("Groups: %+v", groups)
}
