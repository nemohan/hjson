package hjson

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

//TODO: 两种处理错误的策略，一种是遇到错误立即返回
//另一种是记录错误，继续解析, 这里打算采用第一种方法，第二种方法比较适合
//交互的程序编译

//TODO: 优化错误处理
type parser struct {
	jscanner *scanner
	token    int
	literal  string
}

func NewParser(s string) *parser {
	return newParser(bytes.NewBufferString(s))
}
func newParser(reader io.Reader) *parser {
	return &parser{
		jscanner: newScanner(reader),
	}
}
func (p *parser) parse() (Value, error) {
	//获取第一个token
	token, literal := p.jscanner.nextToken()
	p.token = token
	p.literal = literal
	switch token {
	case tokenLBrace:
		return p.parseObject()
	case tokenLBracket:
		return p.parseArray()
	case tokenEOF:
		return nil, errEOF
	case tokenInvalid:
		return nil, p.jscanner.err
	default:
	}
	return nil, fmt.Errorf(`expected: '{', '[' got: %s`, p.literal)
}

func (p *parser) parseObject() (Value, error) {
	obj := NewObject()
	p.match(tokenLBrace)
	for {
		if p.token != tokenString {
			break
		}
		key := p.literal
		if _, ok := obj.values[key]; ok {
			return nil, fmt.Errorf("repeated key:%s in object", key)
		}
		p.match(tokenString)
		if !p.match(tokenColon) {
			return nil, p.getErr(fmt.Errorf(`expect: ':' got:%s`, p.literal))
		}
		value, err := p.parseValues()
		if err != nil {
			return nil, err
		}
		obj.values[key] = value
		if p.token != tokenComma {
			break
		}
		p.match(tokenComma)
	}
	if !p.match(tokenRBrace) {
		return nil, p.getErr(fmt.Errorf("expect: key-value pair or '}' got:%s", p.literal))
	}
	return obj, nil
}

func (p *parser) parseValues() (Value, error) {
	switch p.token {
	case tokenLBrace:
		return p.parseObject()
	case tokenLBracket:
		return p.parseArray()
	case tokenString:
		value := p.literal
		p.match(tokenString)
		return JString(value), nil
	case tokenNumber:
		return p.parseNumber()
	case tokenNull:
		p.match(tokenNull)
		return JNull{}, nil
	case tokenTrue:
		p.match(tokenTrue)
		return JBool(true), nil
	case tokenFalse:
		p.match(tokenFalse)
		return JBool(false), nil
	case tokenEOF:
		return nil, errEOF
	}
	return nil, p.getErr(fmt.Errorf(`expect: STRING, NUMBER, TRUE, FALSE, NULL, {, [`))
}

func (p *parser) parseArray() (Value, error) {
	p.match(tokenLBracket)
	array := NewArray()
	for {
		switch p.token {
		case tokenNumber:
			v, _ := p.parseNumber()
			array.addValue(v)
		case tokenString:
			array.addValue(JString(p.literal))
			p.match(tokenString)
		case tokenNull:
			array.addValue(JNull{})
			p.match(tokenNull)
		case tokenTrue:
			array.addValue(JBool(true))
			p.match(tokenTrue)
		case tokenFalse:
			array.addValue(JBool(false))
			p.match(tokenFalse)
		case tokenLBrace:
			v, err := p.parseObject()
			if err != nil {
				return nil, err
			}
			array.addValue(v)
		case tokenLBracket:
			v, err := p.parseArray()
			if err != nil {
				return nil, err
			}
			array.addValue(v)
		case tokenEOF:
			return nil, errEOF
		}
		if p.token != tokenComma {
			break
		}
		p.match(tokenComma)
	}
	if !p.match(tokenRBracket) {
		return nil, fmt.Errorf("expected:], got: %s", p.literal)
	}
	return array, nil
}

func (p *parser) parseNumber() (Value, error) {
	v, _ := strconv.Atoi(p.literal)
	p.match(tokenNumber)
	return JNumber(v), nil
}
func (p *parser) parseString() (Value, error) {

	return nil, nil
}

func (p *parser) match(token int) bool {
	if p.token == token {
		p.token, p.literal = p.jscanner.nextToken()
		return true
	}
	return false
}

func (p *parser) getErr(err error) error {
	if p.jscanner.err != nil {
		return p.jscanner.err
	}
	return err
}
