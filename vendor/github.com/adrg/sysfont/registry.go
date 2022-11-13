/*
Package sysfont is a small package that makes it easy to identify installed
fonts. It is useful for listing installed fonts or for matching fonts based on
user queries. The matching process also suggests viable font alternatives.
*/
package sysfont

import (
	"path/filepath"
	"strings"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

// Font represents a system font.
type Font struct {
	// Family contains name of the font family.
	Family string

	// Name contains the full name of the font.
	Name string

	// Filename contains the path of the font file.
	Filename string
}

// clone returns a duplicate of the current font instance.
func (f *Font) clone() *Font {
	if f == nil {
		return nil
	}

	font := *f
	return &font
}

type registry struct {
	fonts        []*Font
	families     map[string][]*Font
	filenames    map[string][]*Font
	alternatives [][]string
	defaults     []string
}

func (r *registry) matchFontsByFilename(filename string) []*Font {
	// Attempt to identify font filename in the registry.
	if fonts := r.fontsByFilename(filename); len(fonts) > 0 {
		return fonts
	}

	// Identify font family.
	basename := filepath.Base(filename)
	query := cleanQuery(strings.TrimSuffix(basename, filepath.Ext(basename)))

	queryFamily, ok := r.matchFamily(query)
	if !ok {
		return nil
	}

	// Attempt to identify font by filename and the extracted family.
	match := r.matchFont(query, r.families[queryFamily])
	if match == nil {
		return nil
	}

	// Identify all suitable fonts in the matched family.
	jw := metrics.NewJaroWinkler()
	jw.CaseSensitive = false

	var fonts []*Font
	for _, font := range r.fontsByFilename(match.Filename) {
		if score := (strutil.Similarity(queryFamily, font.Family, jw) +
			strutil.Similarity(query, font.Name, jw)) / 2; score >= 0.97 {
			font.Filename = filename
			fonts = append(fonts, font)
		}
	}

	return fonts
}

func (r *registry) fontsByFilename(filename string) []*Font {
	regFonts, ok := r.filenames[strings.ToLower(filepath.Base(filename))]
	if !ok {
		return nil
	}

	fonts := make([]*Font, len(regFonts))
	for i, regFont := range regFonts {
		fonts[i] = &Font{
			Name:     regFont.Name,
			Family:   regFont.Family,
			Filename: filename,
		}
	}

	return fonts
}

func (r *registry) matchFamily(query string) (string, bool) {
	// Extract font family from query.
	queryFamily := extractFamily(query)

	// Attempt to match extracted family.
	var maxScore float64
	var maxScoreFamily string

	for family := range r.families {
		if score := getFamilyScore(queryFamily, family); score > maxScore {
			maxScore = score
			maxScoreFamily = family
		}
	}

	if maxScore >= 0.97 {
		return maxScoreFamily, true
	}

	return queryFamily, false
}

func (r *registry) matchFont(query string, fonts []*Font) *Font {
	// Clean input and extract font family.
	query = cleanQuery(query)
	queryFamily := extractFamily(query)

	// Attempt to match font.
	var maxScore float64
	var maxScoreFont *Font

	for _, font := range fonts {
		score := getFamilyScore(queryFamily, font.Family)
		if score >= 0.85 {
			score += getFontStyleScore(query, font.Name)
		}
		if score > maxScore {
			maxScore = score
			maxScoreFont = font
		}
	}

	if maxScore >= 0.9 {
		return maxScoreFont
	}

	return nil
}

func (r *registry) getAlternatives(queryFamily string, fonts []*Font) []*Font {
	// Match font family.
	queryFamily = strings.ToLower(queryFamily)

	// Find alternative font families for the extracted family.
	var families []string
	for _, familyGroup := range r.alternatives {
		for _, family := range familyGroup {
			if queryFamily == strings.ToLower(family) {
				families = append(families, familyGroup...)
				break
			}
		}
	}

	// If no alternatives are found, use default families.
	if len(families) == 0 {
		families = r.defaults
	}

	// Match alternative fonts by family.
	var alternatives []*Font
	for _, family := range families {
		family = strings.ToLower(family)
		for _, font := range fonts {
			if family == strings.ToLower(font.Family) {
				alternatives = append(alternatives, font)
			}
		}
	}

	return alternatives
}
