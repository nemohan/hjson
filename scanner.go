package hjson

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"unicode"
)

const (
	tokenInvalid = iota
	tokenComma
	tokenColon
	tokenQuota
	tokenLBrace
	tokenRBrace
	tokenLBracket
	tokenRBracket
	tokenNumber
	tokenString
	tokenNull
	tokenTrue
	tokenFalse
	tokenEOF
)

var (
	tokenTable = map[int]string{
		tokenInvalid:  "invalidToken",
		tokenComma:    ",",
		tokenColon:    ":",
		tokenQuota:    "\"",
		tokenLBrace:   "{",
		tokenRBrace:   "}",
		tokenLBracket: "[",
		tokenRBracket: "]",
		tokenNumber:   "number",
		tokenString:   "string",
		tokenNull:     "null",
		tokenTrue:     "true",
		tokenFalse:    "false",
		tokenEOF:      "eof",
	}
	errEOF = errors.New("unexpected of JSON input")
)

type scanner struct {
	reader *bufio.Reader
	buf    *bytes.Buffer
	line   int
	pos    int
	err    error
}

func newScanner(reader io.Reader) *scanner {
	return &scanner{
		reader: bufio.NewReader(reader),
		buf:    bytes.NewBuffer(make([]byte, 0, 128)),
		line:   1,
		pos:    1,
	}
}

func (s *scanner) nextToken() (int, string) {
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			break
		}
		if isWhitespace(r) {
			continue
		}
		switch r {
		case '"':
			s.buf.Reset()
			if err := s.scanString(); err != nil {
				if err == errEOF {
					goto out
				}
				s.err = err
				return tokenInvalid, s.buf.String()
			}
			return tokenString, s.buf.String()
		case ',':
			return tokenComma, ","
		case ':':
			return tokenColon, ":"
		case '{':
			return tokenLBrace, "{"
		case '}':
			return tokenRBrace, "}"
		case '[':
			return tokenLBracket, "["
		case ']':
			return tokenRBracket, "]"
		case '\n':
			s.line++
		default:
			s.buf.Reset()
			s.reader.UnreadRune()
			if unicode.IsDigit(r) {
				if err := s.scanNumber(); err != nil {
					if err == errEOF {
						goto out
					}
					s.err = err
					return tokenInvalid, s.buf.String()
				}
				return tokenNumber, s.buf.String()
			} else {
				if err := s.scanIdent(); err != nil {
					goto out
				}
				lit := s.buf.String()
				return lookup(lit), lit
			}
		}
	}
out:
	return tokenEOF, ""
}

func lookup(lit string) int {
	switch lit {
	case "null":
		return tokenNull
	case "true":
		return tokenTrue
	case "false":
		return tokenFalse
	}
	return tokenInvalid
}

func (s *scanner) scanHex() error {
	for i := 0; i < 4; i++ {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			return errEOF
		}
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F') {
			s.buf.WriteRune(r)
			continue
		}
		return fmt.Errorf("encounter invalid hexadecimal digit:%s", string(r))
	}
	return nil
}
func (s *scanner) scanString() error {
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			//return errEOF
			break
		}
		if r == '\\' {
			s.buf.WriteRune(r)
			r, _, err = s.reader.ReadRune()
			if err != nil {
				return errEOF
			}
			switch r {
			case '"':
			case '\\':
			case '/':
			case 'b':
			case 'f':
			case 'n':
			case 'r':
			case 't':
			case 'u':
				s.buf.WriteRune(r)
				if err := s.scanHex(); err != nil {
					return err
				}
				goto next
			default:
				return fmt.Errorf("invalid escape sequence: %s", string(r))
			}
		} else if r == '"' {
			break
		}
		s.buf.WriteRune(r)
	next:
	}
	return nil
}

//.12 or .12e13 .12E13 .12e+13 .12e-13
func (s *scanner) scanFraction() error {
	sawE := false
	sawSign := false
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			return errEOF
		}
		switch r {
		case '-', '+':
			if !sawSign {
				s.buf.WriteRune(r)
				sawSign = true
			} else {
				return fmt.Errorf(`invalid character:'%s' in numeric literal`, string(r))
			}
		case 'e', 'E':
			if sawSign && !sawE {
				return fmt.Errorf(`invalid character:'%s' in numeric literal`, string(r))
			} else if !sawE {
				s.buf.WriteRune(r)
				sawE = true
			} else {
				return fmt.Errorf(`invalid character:'%s' in numeric literal`, string(r))
			}
		default:
			if unicode.IsDigit(r) {
				s.buf.WriteRune(r)
			} else {
				s.reader.UnreadRune()
				return nil
			}
		}
	}
}

func (s *scanner) scanNumber() error {
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsDigit(r) {
			s.buf.WriteRune(r)
		} else if r == '.' { //int frac
			s.buf.WriteRune(r)
			return s.scanFraction()
		} else if r == 'e' || r == 'E' { //int exp
			s.buf.WriteRune(r)
			return s.scanFraction()
		} else {
			if !unicode.IsDigit(r) {
				s.reader.UnreadRune()
				break
			}
		}
	}
	return nil
}
func (s *scanner) scanIdent() error {
	for {
		r, _, err := s.reader.ReadRune()
		if err != nil {
			return errEOF
		}
		if unicode.IsLetter(r) {
			s.buf.WriteRune(r)
		} else {
			s.reader.UnreadRune()
			break
		}
	}
	return nil
}
func isWhitespace(r rune) bool {
	if r == ' ' || r == '\t' || r == '\r' {
		return true
	}
	return false
}
