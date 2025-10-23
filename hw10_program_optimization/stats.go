package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

// функцию не трогаем.
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	// используем bufio и easyjson
	var i int
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var user User
		if err = user.UnmarshalJSON(sc.Bytes()); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	// предварительная компиляция регулярки
	compile, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		if compile.Match([]byte(user.Email)) {
			userSplit := strings.SplitN(user.Email, "@", 2)
			domainName := strings.ToLower(userSplit[1])
			result[domainName]++
		}
	}

	return result, nil
}
