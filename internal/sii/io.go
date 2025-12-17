package sii

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// ReadDocument handles plaintext SII content:
//   - data doit déjà être déchiffré (TODO: appel à un outil externe pour les fichiers chiffrés)
//   - parse le texte en structure Document.
func ReadDocument(data []byte) (*Document, error) {
	// Si le fichier commence directement par "SiiNunit", on considère que c'est du texte clair.
	if !bytes.Contains(data, []byte("SiiNunit")) {
		// Ici il faudrait appeler un outil natif (sii_decrypt / SII_Decrypt.dll).
		return nil, errors.New("encrypted or unknown SII format: decryption support not yet implemented")
	}

	lines := splitLines(data)

	if len(lines) < 2 || !strings.HasPrefix(strings.TrimSpace(lines[0]), "SiiNunit") {
		return nil, errors.New("invalid SII: missing SiiNunit header")
	}

	// Chercher la première accolade ouvrante suivant SiiNunit
	i := 1
	for i < len(lines) && strings.TrimSpace(lines[i]) != "{" {
		i++
	}
	if i == len(lines) {
		return nil, errors.New("invalid SII: missing opening brace")
	}
	i++ // première ligne après "{"

	var blocks []Block

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" || line == "}" {
			i++
			continue
		}

		// Ligne de bloc: "type : name {"
		if strings.Contains(line, ":") && strings.Contains(line, "{") {
			parts := strings.SplitN(line, ":", 2)
			typ := strings.TrimSpace(parts[0])

			right := parts[1]
			right = strings.SplitN(right, "{", 2)[0]
			name := strings.TrimSpace(right)

			// Collecter les lignes jusqu'à la prochaine "}"
			i++
			var body []string
			for i < len(lines) {
				l := lines[i]
				body = append(body, l)
				if strings.TrimSpace(l) == "}" {
					i++
					break
				}
				i++
			}

			blocks = append(blocks, Block{
				Type:       typ,
				Name:       name,
				Properties: parseProperties(body),
			})
			continue
		}

		i++
	}

	return &Document{Blocks: blocks}, nil
}

// WriteDocument sérialise un Document vers du texte SII. La partie chiffrement
// éventuelle doit être appliquée par l'appelant si nécessaire.
func WriteDocument(doc *Document) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("SiiNunit\n{\n")

	for _, b := range doc.Blocks {
		fmt.Fprintf(&buf, "%s : %s {\n", b.Type, b.Name)
		for k, vals := range b.Properties {
			for _, v := range vals {
				fmt.Fprintf(&buf, " %s: %s\n", k, v)
			}
		}
		buf.WriteString("}\n")
	}

	buf.WriteString("}\n")
	return buf.Bytes(), nil
}

func splitLines(b []byte) []string {
	sc := bufio.NewScanner(bytes.NewReader(b))
	var out []string
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out
}

func parseProperties(lines []string) map[string][]string {
	props := make(map[string][]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || l == "}" || l == "{" {
			continue
		}
		parts := strings.SplitN(l, ":", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		props[k] = append(props[k], v)
	}
	return props
}
