package main

import (
    "fmt"
    "regexp"
    "strings"
)

const ID_FORMAT = "[a-z]+[0-9]+[a-z0-9]+"

var allowedList = [...]*regexp.Regexp{
    regexp.MustCompile(`^company$`),                                                    // /company/
    regexp.MustCompile(strings.Replace(`^company/{id}$`, "{id}", ID_FORMAT, 1)),        // /company/{id}
    regexp.MustCompile(strings.Replace(`^company/[a-z0-9]+$`, "{id}", ID_FORMAT, 1)),   // /company/{id}
    regexp.MustCompile(`^company/account$`),                                            // /company/account
    regexp.MustCompile(`^account$`),                                                    // /account
    regexp.MustCompile(strings.Replace(`^account/{id}$`, "{id}", ID_FORMAT, 1)),        // /account/{id}
    regexp.MustCompile(strings.Replace(`^{id}$`, "{id}", ID_FORMAT, 1)),                // /{id}
    regexp.MustCompile(strings.Replace(`^account/{id}/user$`, "{id}", ID_FORMAT, 1)),   // /?account/{id}/user
    regexp.MustCompile(`^tenant/account/blocked$`),                                     // /tenant/account/blocked
}

func main() {
    fmt.Println("Running ...")
}

func ValidatePath(path string) bool {
    for _, allowedPath := range allowedList {
        if allowedPath.MatchString(strings.Trim(strings.ToLower(path), "/")) {
            return true
        }
    }

    return false
}
