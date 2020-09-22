package main

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestValidator(t *testing.T) {
    var cases = []struct {
        path      string
        expection bool
    }{
        {"company", true},
        {"tenant/sj3co3s4", false},
        {"company/sd45f768", true},
        {"account/acc74850", true},
        {"company/account", true},
        {"acc734340", true},
        {"account/acc234234/user", true},
        {"account/blocked", false},
        {"tenant/account/blocked", true},
        {"tenant/account/acc23849", false},
    }

    proxy := Proxy{
        allowedList: allowedList,
    }

    for _, tc := range cases {
        require.Equal(t, tc.expection, proxy.ValidatePath(tc.path), "Test is failing! %d equal %s", tc.expection, tc.path)
    }
}

