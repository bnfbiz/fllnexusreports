package csvmap

// Modified version of package created by CoPilot to allow for headers to be dynamicely read from the CSV file.
import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// Normalize a header to lowercase, trimmed and alphanumeric-only.
func normalizeHeader(s string) string {
	t := strings.ToLower(strings.TrimSpace(s))
	b := make([]rune, 0, len(t))
	for _, r := range t {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b = append(b, r)
		}
	}
	return string(b)
}

// ReadCSVToMaps reads a CSV file and returns a slice of maps[string]string
// keyed by canonicalName -> string. The aliases parameter maps a canonical
// name to a list of possible header aliases for that field in CSV headers.
func ReadCSVToMaps(path string) ([]map[string]string, map[string]int, map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // accept variable-length rows

	header, err := r.Read()
	if err == io.EOF {
		return nil, nil, nil, fmt.Errorf("csv has no header")
	}
	if err != nil {
		return nil, nil, nil, err
	}

	// build header lookup: normalized header string -> index
	headMap := map[string]int{}
	headerNames := map[string]string{}
	for i, h := range header {
		normalized := normalizeHeader(h)
		headMap[normalized] = i
		headerNames[normalized] = h
	}

	out := []map[string]string{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, err
		}
		m := map[string]string{}
		for canonical, idx := range headMap {
			if idx >= 0 && idx < len(rec) {
				m[canonical] = rec[idx]
			} else {
				m[canonical] = ""
			}
		}
		out = append(out, m)
	}

	keys := []string{}
	for k := range headMap {
		keys = append(keys, k)
	}

	return out, headMap, headerNames, nil
}
