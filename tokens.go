package main

import (
	"go/token"
	"strconv"
)

var tokens = [...]string{
	token.ILLEGAL:        "ILLEGAL",
	token.EOF:            "EOF",
	token.COMMENT:        "COMMENT",
	token.IDENT:          "IDENT",
	token.INT:            "INT",
	token.FLOAT:          "FLOAT",
	token.IMAG:           "IMAG",
	token.CHAR:           "CHAR",
	token.STRING:         "STRING",
	token.ADD:            "ADD",
	token.SUB:            "SUB",
	token.MUL:            "MUL",
	token.QUO:            "QUO",
	token.REM:            "REM",
	token.AND:            "AND",
	token.OR:             "OR",
	token.XOR:            "XOR",
	token.SHL:            "SHL",
	token.SHR:            "SHR",
	token.AND_NOT:        "AND_NOT",
	token.ADD_ASSIGN:     "ADD_ASSIGN",
	token.SUB_ASSIGN:     "SUB_ASSIGN",
	token.MUL_ASSIGN:     "MUL_ASSIGN",
	token.QUO_ASSIGN:     "QUO_ASSIGN",
	token.REM_ASSIGN:     "REM_ASSIGN",
	token.AND_ASSIGN:     "AND_ASSIGN",
	token.OR_ASSIGN:      "OR_ASSIGN",
	token.XOR_ASSIGN:     "XOR_ASSIGN",
	token.SHL_ASSIGN:     "SHL_ASSIGN",
	token.SHR_ASSIGN:     "SHR_ASSIGN",
	token.AND_NOT_ASSIGN: "AND_NOT_ASSIGN",
	token.LAND:           "LAND",
	token.LOR:            "LOR",
	token.ARROW:          "ARROW",
	token.INC:            "INC",
	token.DEC:            "DEC",
	token.EQL:            "EQL",
	token.LSS:            "LSS",
	token.GTR:            "GTR",
	token.ASSIGN:         "ASSIGN",
	token.NOT:            "NOT",
	token.NEQ:            "NEQ",
	token.LEQ:            "LEQ",
	token.GEQ:            "GEQ",
	token.DEFINE:         "DEFINE",
	token.ELLIPSIS:       "ELLIPSIS",
	token.LPAREN:         "LPAREN",
	token.LBRACK:         "LBRACK",
	token.LBRACE:         "LBRACE",
	token.COMMA:          "COMMA",
	token.PERIOD:         "PERIOD",
	token.RPAREN:         "RPAREN",
	token.RBRACK:         "RBRACK",
	token.RBRACE:         "RBRACE",
	token.SEMICOLON:      "SEMICOLON",
	token.COLON:          "COLON",
	token.BREAK:          "BREAK",
	token.CASE:           "CASE",
	token.CHAN:           "CHAN",
	token.CONST:          "CONST",
	token.CONTINUE:       "CONTINUE",
	token.DEFAULT:        "DEFAULT",
	token.DEFER:          "DEFER",
	token.ELSE:           "ELSE",
	token.FALLTHROUGH:    "FALLTHROUGH",
	token.FOR:            "FOR",
	token.FUNC:           "FUNC",
	token.GO:             "GO",
	token.GOTO:           "GOTO",
	token.IF:             "IF",
	token.IMPORT:         "IMPORT",
	token.INTERFACE:      "INTERFACE",
	token.MAP:            "MAP",
	token.PACKAGE:        "PACKAGE",
	token.RANGE:          "RANGE",
	token.RETURN:         "RETURN",
	token.SELECT:         "SELECT",
	token.STRUCT:         "STRUCT",
	token.SWITCH:         "SWITCH",
	token.TYPE:           "TYPE",
	token.VAR:            "VAR",
}

func tokenString(tok token.Token) string {
	s := ""
	if 0 <= tok && tok < token.Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}
