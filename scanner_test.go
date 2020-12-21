package hjson

import (
	"bytes"
	"testing"
)

func TestScanner(t *testing.T) {
	type token struct {
		tokenType int
		literal   string
	}
	testCases := []struct {
		json   string
		tokens []token
	}{
		{
			json: `[123, "token"]`,
			tokens: []token{
				{tokenLBracket, "["},
				{tokenNumber, "123"},
				{tokenComma, ","},
				{tokenString, "token"},
				{tokenRBracket, "]"},
				{tokenEOF, ""},
			},
		},
		{
			json: `{"xx":"k\"ey", "value":[1, 2, true, null]}`,
			tokens: []token{
				{tokenLBrace, "{"},
				{tokenString, "xx"},
				{tokenColon, ":"},
				{tokenString, `k\"ey`},
				{tokenComma, ","},
				{tokenString, "value"},
				{tokenColon, ":"},
				{tokenLBracket, "["},
				{tokenNumber, "1"},
				{tokenComma, ","},
				{tokenNumber, "2"},
				{tokenComma, ","},
				{tokenTrue, "true"},
				{tokenComma, ","},
				{tokenNull, "null"},
				{tokenRBracket, "]"},
				{tokenRBrace, "}"},
				{tokenEOF, ""},
			},
		},
		{
			json: `{"key":"\"\\ha\/\b\f\n\r\t"}`,
			tokens: []token{
				{tokenLBrace, "{"},
				{tokenString, "key"},
				{tokenColon, ":"},
				{tokenString, `\"\\ha\/\b\f\n\r\t`},
				{tokenRBrace, "}"},
			},
		},
		{
			json: `{"key":"\"hh"}`,
			tokens: []token{
				{tokenLBrace, "{"},
				{tokenString, "key"},
				{tokenColon, ":"},
				{tokenString, `\"hh`},
				{tokenRBrace, "}"},
			},
		},
		{
			json: `[12,12.12, 12e+1, 12e-1,12.12e+1]`,
			tokens: []token{
				{tokenLBracket, "["},
				{tokenNumber, "12"},
				{tokenComma, ","},

				{tokenNumber, "12.12"},
				{tokenComma, ","},

				{tokenNumber, "12e+1"},
				{tokenComma, ","},

				{tokenNumber, "12e-1"},
				{tokenComma, ","},

				{tokenNumber, "12.12e+1"},
				{tokenRBracket, "]"},
				{tokenEOF, ""},
			},
		},
	}

	for i, tc := range testCases {
		scanner := newScanner(bytes.NewBufferString(tc.json))
		for _, token := range tc.tokens {
			tok, literal := scanner.nextToken()
			if token.tokenType != tok || token.literal != literal {
				t.Fatalf("case:%d expect:<%s %s> got:<%s %s>\n", i, tokenTable[token.tokenType],
					token.literal, tokenTable[tok], literal)
			}
		}
	}

}
