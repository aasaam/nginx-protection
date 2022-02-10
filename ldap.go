package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func ldapLogin(
	url string,
	readOnlyUsername string,
	readOnlyPassword string,
	baseDN string,
	filter string,
	attributesJSON string,
	username string,
	password string,
) (bool, error, error) {

	l, connectionError := ldap.DialURL(url, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if connectionError != nil {
		return false, nil, errors.New("ldap connection failed: " + connectionError.Error())
	}
	defer l.Close()

	readonlyBindError := l.Bind(readOnlyUsername, readOnlyPassword)
	if readonlyBindError != nil {
		return false, nil, errors.New("ldap readonly user bind failed: " + readonlyBindError.Error())
	}

	filter = strings.ReplaceAll(filter, ldapUsernamePlaceHolder, username)

	var attributes []string
	jsonError := json.Unmarshal([]byte(attributesJSON), &attributes)
	if jsonError != nil {
		return false, nil, errors.New("ldap configure json attributes failed: " + jsonError.Error())
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)

	searchResult, err := l.Search(searchRequest)
	if err != nil {
		return false, nil, errors.New("ldap search failed: " + jsonError.Error())
	}

	count := len(searchResult.Entries)

	if count > 1 {
		return false, nil, errors.New("to many user returns, check your filter and ldap properties")
	}

	if count <= 0 {
		return false, errors.New("user does not exist"), nil
	}

	err = l.Bind(searchResult.Entries[0].DN, password)
	if err != nil {
		return false, errors.New("user dn not found or password is incorrect"), nil
	}

	return true, nil, nil
}
