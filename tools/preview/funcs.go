package preview

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
)

// TemplateFuncs maps Gitea's built-in template functions available in mail templates.
var TemplateFuncs = template.FuncMap{
	"AppUrl":             func() string { return "https://gitea.com/" },
	"AppName":            func() string { return "Gitea" },
	"DotEscape":          func(s string) template.HTML { return template.HTML(template.HTMLEscapeString(s)) },
	"QueryEscape":        url.QueryEscape,
	"ShortSha":           ShortSha,
	"PathEscapeSegments": PathEscapeSegments,
	"FileSize":           FileSize,
	"HTMLFormat":         HTMLFormat,
	"Dict":               Dict,
	"gt":                 Gt,  // override built-in for JSON float64/int coercion
	"lt":                 Lt,
	"ge":                 Ge,
	"le":                 Le,
}

// ShortSha truncates a commit SHA to 10 characters.
func ShortSha(s string) string {
	if len(s) > 10 {
		return s[:10]
	}
	return s
}

// PathEscapeSegments splits a path and URL-escapes each segment.
func PathEscapeSegments(s string) string {
	parts := strings.Split(s, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}

// FileSize formats a byte count as a human-readable string (e.g. "2.5 MB").
func FileSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for nn := n / unit; nn >= unit && exp < 3; nn /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %s", float64(n)/float64(div), []string{"KB", "MB", "GB", "TB"}[exp])
}

// HTMLFormat substitutes %s/%d placeholders with HTML-escaped values.
func HTMLFormat(fmtStr string, args ...any) template.HTML {
	result := fmtStr
	for i, arg := range args {
		escaped := template.HTMLEscapeString(fmt.Sprint(arg))
		result = strings.Replace(result, fmt.Sprintf("%%[%d]s", i+1), escaped, 1)
		if !strings.Contains(result, fmt.Sprintf("%%[%d]s", i+1)) {
			result = strings.Replace(result, "%s", escaped, 1)
			result = strings.Replace(result, "%d", fmt.Sprint(arg), 1)
		}
	}
	return template.HTML(result)
}

// Dict builds a map from alternating key/value pairs (used by Gitea's Dict function).
func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("odd number of Dict args")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("Dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// toFloat64 coerces a value to float64 if numeric.
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case int16:
		return float64(val), true
	case int8:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint64:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint8:
		return float64(val), true
	default:
		return 0, false
	}
}

// Gt is a type-coercing > comparison that handles JSON float64 vs int literals.
func Gt(a, b any) bool {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af > bf
	}
	// Fallback to string comparison
	return fmt.Sprint(a) > fmt.Sprint(b)
}

// Lt is a type-coercing < comparison that handles JSON float64 vs int literals.
func Lt(a, b any) bool {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af < bf
	}
	return fmt.Sprint(a) < fmt.Sprint(b)
}

// Ge is a type-coercing >= comparison.
func Ge(a, b any) bool {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af >= bf
	}
	return fmt.Sprint(a) >= fmt.Sprint(b)
}

// Le is a type-coercing <= comparison.
func Le(a, b any) bool {
	af, aok := toFloat64(a)
	bf, bok := toFloat64(b)
	if aok && bok {
		return af <= bf
	}
	return fmt.Sprint(a) <= fmt.Sprint(b)
}
