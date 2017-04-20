package inflector

import (
	"regexp"
	"strings"
)

type (
	rule struct {
		re   *regexp.Regexp
		repl string
	}
)

var (
	// https://github.com/cakephp/cakephp/blob/master/src/Utility/Inflector.php#L33-L57
	pluralRules = []rule{
		{re: regexp.MustCompile(`(?i)(s)tatus$`), repl: "${1}tatuses"},
		{re: regexp.MustCompile(`(?i)(quiz)$`), repl: "${1}zes"},
		{re: regexp.MustCompile(`(?i)^(ox)$`), repl: "${1}${2}en"},
		{re: regexp.MustCompile(`(?i)([m|l])ouse$`), repl: "${1}ice"},
		{re: regexp.MustCompile(`(?i)(matr|vert|ind)(ix|ex)$`), repl: "${1}ices"},
		{re: regexp.MustCompile(`(?i)(x|ch|ss|sh)$`), repl: "${1}es"},
		{re: regexp.MustCompile(`(?i)([^aeiouy]|qu)y$`), repl: "${1}ies"},
		{re: regexp.MustCompile(`(?i)(hive)$`), repl: "${1}s"},
		{re: regexp.MustCompile(`(?i)(chef)$`), repl: "${1}s"},
		{re: regexp.MustCompile(`(?i)(?:([^f])fe|([lre])f)$`), repl: "${1}${2}ves"},
		{re: regexp.MustCompile(`(?i)sis$`), repl: "ses"},
		{re: regexp.MustCompile(`(?i)([ti])um$`), repl: "${1}a"},
		{re: regexp.MustCompile(`(?i)(p)erson$`), repl: "${1}eople"},
		//rule{re: regexp.MustCompile(`(?i)(?<!u)(m)an$`), repl: "${1}en"}, // TODO regexp compile error
		{re: regexp.MustCompile(`(?i)(m)an$`), repl: "${1}en"},
		{re: regexp.MustCompile(`(?i)(c)hild$`), repl: "${1}hildren"},
		{re: regexp.MustCompile(`(?i)(buffal|tomat)o$`), repl: "${1}${2}oes"},
		{re: regexp.MustCompile(`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin)us$`), repl: "${1}i"},
		{re: regexp.MustCompile(`(?i)us$`), repl: "uses"},
		{re: regexp.MustCompile(`(?i)(alias)$`), repl: "${1}es"},
		{re: regexp.MustCompile(`(?i)(ax|cris|test)is$`), repl: "${1}es"},
		{re: regexp.MustCompile(`s$`), repl: "s"},
		{re: regexp.MustCompile(`^$`), repl: ""},
		{re: regexp.MustCompile(`$`), repl: "s"},
	}
	// https://github.com/cakephp/cakephp/blob/master/src/Utility/Inflector.php#L64-L99
	singularRules = []rule{
		{re: regexp.MustCompile(`(?i)(s)tatuses$`), repl: "${1}${2}tatus"},
		{re: regexp.MustCompile(`(?i)^(.*)(menu)s$`), repl: "${1}${2}"},
		{re: regexp.MustCompile(`(?i)(quiz)zes$`), repl: "\\1"},
		{re: regexp.MustCompile(`(?i)(matr)ices$`), repl: "${1}ix"},
		{re: regexp.MustCompile(`(?i)(vert|ind)ices$`), repl: "${1}ex"},
		{re: regexp.MustCompile(`(?i)^(ox)en`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(alias)(es)*$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$`), repl: "${1}us"},
		{re: regexp.MustCompile(`(?i)([ftw]ax)es$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(cris|ax|test)es$`), repl: "${1}is"},
		{re: regexp.MustCompile(`(?i)(shoe)s$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(o)es$`), repl: "${1}"},
		{re: regexp.MustCompile(`ouses$`), repl: "ouse"},
		{re: regexp.MustCompile(`([^a])uses$`), repl: "${1}us"},
		{re: regexp.MustCompile(`(?i)([m|l])ice$`), repl: "${1}ouse"},
		{re: regexp.MustCompile(`(?i)(x|ch|ss|sh)es$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(m)ovies$`), repl: "${1}${2}ovie"},
		{re: regexp.MustCompile(`(?i)(s)eries$`), repl: "${1}${2}eries"},
		{re: regexp.MustCompile(`(?i)([^aeiouy]|qu)ies$`), repl: "${1}y"},
		{re: regexp.MustCompile(`(?i)(tive)s$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(hive)s$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)(drive)s$`), repl: "${1}"},
		{re: regexp.MustCompile(`(?i)([le])ves$`), repl: "${1}f"},
		{re: regexp.MustCompile(`(?i)([^rfoa])ves$`), repl: "${1}fe"},
		{re: regexp.MustCompile(`(?i)(^analy)ses$`), repl: "${1}sis"},
		{re: regexp.MustCompile(`(?i)(analy|diagno|^ba|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$`), repl: "${1}${2}sis"},
		{re: regexp.MustCompile(`(?i)([ti])a$`), repl: "${1}um"},
		{re: regexp.MustCompile(`(?i)(p)eople$`), repl: "${1}${2}erson"},
		{re: regexp.MustCompile(`(?i)(m)en$`), repl: "${1}an"},
		{re: regexp.MustCompile(`(?i)(c)hildren$`), repl: "${1}${2}hild"},
		{re: regexp.MustCompile(`(?i)(n)ews$`), repl: "${1}${2}ews"},
		{re: regexp.MustCompile(`eaus$`), repl: "eau"},
		{re: regexp.MustCompile(`^(.*us)$`), repl: "\\1"},
		{re: regexp.MustCompile(`(?i)s$`), repl: ""},
	}
	irregularToPluralMap = map[string]string{
		"atlas":     "atlases",
		"beef":      "beefs",
		"brief":     "briefs",
		"brother":   "brothers",
		"cafe":      "cafes",
		"child":     "children",
		"cookie":    "cookies",
		"corpus":    "corpuses",
		"cow":       "cows",
		"criterion": "criteria",
		"ganglion":  "ganglions",
		"genie":     "genies",
		"genus":     "genera",
		"graffito":  "graffiti",
		"hoof":      "hoofs",
		"loaf":      "loaves",
		"man":       "men",
		"money":     "monies",
		"mongoose":  "mongooses",
		"move":      "moves",
		"mythos":    "mythoi",
		"niche":     "niches",
		"numen":     "numina",
		"occiput":   "occiputs",
		"octopus":   "octopuses",
		"opus":      "opuses",
		"ox":        "oxen",
		"penis":     "penises",
		"person":    "people",
		"sex":       "sexes",
		"soliloquy": "soliloquies",
		"testis":    "testes",
		"trilby":    "trilbys",
		"turf":      "turfs",
		"potato":    "potatoes",
		"hero":      "heroes",
		"tooth":     "teeth",
		"goose":     "geese",
		"foot":      "feet",
		"foe":       "foes",
		"sieve":     "sieves",
	}
	irregularToSingularMap = reverseMapKeyValue(irregularToPluralMap)
)

func reverseMapKeyValue(mp map[string]string) map[string]string {
	res := map[string]string{}
	for key, value := range mp {
		res[value] = key
	}
	return res
}

func Pluralize(word string) string {
	if word == "" {
		return ""
	}
	if v := irregularToPluralMap[word]; v != "" {
		return v
	}
	for _, rule := range pluralRules {
		if rule.re.MatchString(word) {
			return rule.re.ReplaceAllString(word, rule.repl)
		}
	}
	return word + "s" // TODO
}

func Singularize(word string) string {
	if word == "" {
		return ""
	}
	if v := irregularToSingularMap[word]; v != "" {
		return v
	}
	for _, rule := range singularRules {
		if rule.re.MatchString(word) {
			return rule.re.ReplaceAllString(word, rule.repl)
		}
	}
	// TODO
	return strings.TrimSuffix(word, "s")
}
