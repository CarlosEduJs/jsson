package lsp

import (
	"encoding/json"
	"strings"
)

func (s *Server) handleSemanticTokensFull(id, raw json.RawMessage) error {
	var params SemanticTokensParams
	if err := json.Unmarshal(raw, &params); err != nil {
		return s.sendError(id, "Invalid params")
	}

	doc, ok := s.getDocument(params.TextDocument.URI)
	if !ok {
		return s.sendError(id, "Document not found")
	}

	tokens := s.generateSemanticTokens(doc.Content)

	response := ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result: SemanticTokens{
			Data: tokens,
		},
	}

	return s.writeMessage(response)
}

const (
	TokenTypeNamespace     = 0
	TokenTypeType          = 1
	TokenTypeClass         = 2
	TokenTypeEnum          = 3
	TokenTypeInterface     = 4
	TokenTypeStruct        = 5
	TokenTypeTypeParameter = 6
	TokenTypeParameter     = 7
	TokenTypeVariable      = 8
	TokenTypeProperty      = 9
	TokenTypeEnumMember    = 10
	TokenTypeEvent         = 11
	TokenTypeFunction      = 12
	TokenTypeMethod        = 13
	TokenTypeMacro         = 14
	TokenTypeKeyword       = 15
	TokenTypeModifier      = 16
	TokenTypeComment       = 17
	TokenTypeString        = 18
	TokenTypeNumber        = 19
	TokenTypeRegexp        = 20
	TokenTypeOperator      = 21
)

func (s *Server) generateSemanticTokens(content string) []int {
	tokens := []int{}
	lines := strings.Split(content, "\n")

	declaredVars := make(map[string]bool)
	scopeParams := make(map[string]bool)

	prevLine := 0
	prevChar := 0

	for lineNum, line := range lines {
		if before, _, ok := strings.Cut(line, ":="); ok {
			beforeAssign := before

			varName := strings.TrimSpace(beforeAssign)
			if isValidIdentifier(varName) {
				declaredVars[varName] = true

				tokens = append(tokens,
					lineNum-prevLine,
					len(line)-len(strings.TrimLeft(line, " \t"))-prevChar,
					len(varName),
					TokenTypeVariable,
					0,
				)
				prevLine = lineNum
				prevChar = len(line) - len(strings.TrimLeft(line, " \t"))
			}
		}

		if strings.Contains(line, " map (") || strings.Contains(line, " zip (") {
			startIdx := strings.Index(line, " map (")
			if startIdx == -1 {
				startIdx = strings.Index(line, " zip (")
			}

			if startIdx != -1 {
				rest := line[startIdx:]
				parenStart := strings.Index(rest, "(")

				parenEnd := strings.Index(rest, ")")
				if parenStart != -1 && parenEnd != -1 {
					paramsStr := rest[parenStart+1 : parenEnd]
					for _, param := range strings.Split(paramsStr, ",") {
						param = strings.TrimSpace(param)
						if param != "" && isValidIdentifier(param) {
							scopeParams[param] = true

							paramPos := strings.Index(line, param)
							if paramPos != -1 {
								tokens = append(tokens,
									lineNum-prevLine,
									paramPos-prevChar,
									len(param),
									TokenTypeParameter,
									0,
								)
								prevLine = lineNum
								prevChar = paramPos
							}
						}
					}
				}
			}
		}

		for param := range scopeParams {
			idx := 0

			for {
				pos := strings.Index(line[idx:], param)
				if pos == -1 {
					break
				}

				actualPos := idx + pos
				if actualPos > 0 && isIdentifierChar(line[actualPos-1]) {
					idx = actualPos + 1

					continue
				}

				if actualPos+len(param) < len(line) && isIdentifierChar(line[actualPos+len(param)]) {
					idx = actualPos + 1

					continue
				}

				tokens = append(tokens,
					lineNum-prevLine,
					actualPos-prevChar,
					len(param),
					TokenTypeParameter,
					0,
				)
				prevLine = lineNum
				prevChar = actualPos
				idx = actualPos + len(param)
			}
		}
	}

	return tokens
}
