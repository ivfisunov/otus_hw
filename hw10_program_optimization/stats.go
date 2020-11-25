package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User
var user User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	i := 0
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			mail := strings.Split(user.Email, "@")
			splitedMail := strings.ToLower(mail[1])
			result[splitedMail]++
		}
	}
	return result, nil
}
