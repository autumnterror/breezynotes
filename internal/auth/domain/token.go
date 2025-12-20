package domain

import (
	brzrpc "github.com/autumnterror/breezynotes/api/proto/gen"
	"time"
)

type Token struct {
	Value string    `json:"value"`
	Exp   time.Time `json:"exp"`
}

type Tokens struct {
	Access  Token `json:"access"`
	Refresh Token `json:"refresh"`
}

func TokenToRPC(t *Token) *brzrpc.Token {
	if t == nil {
		return nil
	}
	return &brzrpc.Token{
		Value: t.Value,
		Exp:   t.Exp.Unix(),
	}
}

func TokensToRPC(ts *Tokens) *brzrpc.Tokens {
	if ts == nil {
		return nil
	}
	return &brzrpc.Tokens{
		AccessToken:  ts.Access.Value,
		RefreshToken: ts.Refresh.Value,
		ExpAccess:    ts.Access.Exp.Unix(),
		ExpRefresh:   ts.Refresh.Exp.Unix(),
	}
}
