// THIS FILE IS AUTOMATICALLY GENERATED BY `internal/cmd/test_gen` DO NOT EDIT

package cook

import (
	"fmt"
	"reflect"
	"testing"
)

// --------------------------------------------------------------
// Canonical Unit Tests
//
// NOTE: I did ignore `Qty: "some"` for one word ingredients and `Qty: 1` for
// one word cookware. It's trivial to add, however, it is a little weird for
// some ingredients and some cookware i.e. "some egg" or "1 tongs".
//
// It's also presumptious of usage for non-english languages.
// I'm going to leave it to the projects interfacing this parser to work that
// out as they see fit.
//
// Tests are defined here:
// https://github.com/cooklang/spec/tree/fa9bc51515b3317da434cb2b5a4a6ac12257e60b/tests
// --------------------------------------------------------------

// --------------------------------------------------------------
// Utilities
// --------------------------------------------------------------

// Deep compares two recipes after adjusting the format to match the .yaml
// format.
//
// Specifically, this forces pre-parsed fractions in Qty, just like the canonical tests
//
//	e.g.
//	src = "@carrots{1/2}"
//	results in Qty being 0.5 in the canonical tests (instead of 1/2 with QtyVal=0.5)
//
// I do wish I could avoid doing this, but interfaces in go are extremely unwieldy
// so I feel this is likely the better approach.
//
// This also inserts `Some` and 1 for cookware to match tests which I did
// not agree with doing during parsing. See the above `NOTE` in this file.
func assertCanonicalRecipe(t *testing.T, got *Recipe, want *Recipe) {
	// Push the value into Qty as a string
	for i, ingr := range got.Ingredients {
		if ingr.QtyVal != NoQty {
			got.Ingredients[i].Qty = fmt.Sprintf("%v", ingr.QtyVal)
		} else if ingr.Qty == "" {
			// Include the contentious `some`
			got.Ingredients[i].Qty = "some"
		}
	}
	for i, ware := range got.Cookware {
		if ware.QtyVal != NoQty {
			got.Cookware[i].Qty = fmt.Sprintf("%v", ware.QtyVal)
		} else if ware.Qty == "" {
			// Include the contentious `1`
			got.Cookware[i].Qty = "1"
			got.Cookware[i].QtyVal = 1
		}
	}
	for i, timer := range got.Timers {
		if timer.QtyVal != NoQty {
			got.Timers[i].Qty = fmt.Sprintf("%v", timer.QtyVal)
		}
	}

	// Repeat the above for Steps
	for i, step := range got.Steps {
		for j, rawChunk := range step {
			var newChunk Chunk
			switch chunk := rawChunk.(type) {
			case Text:
				newChunk = chunk
			case Ingredient:
				if chunk.QtyVal != NoQty {
					chunk.Qty = fmt.Sprintf("%v", chunk.QtyVal)
				} else if chunk.Qty == "" {
					chunk.Qty = "some"
				}
				newChunk = chunk
			case Cookware:
				if chunk.QtyVal != NoQty {
					chunk.Qty = fmt.Sprintf("%v", chunk.QtyVal)
				} else if chunk.Qty == "" {
					chunk.Qty = "1"
					chunk.QtyVal = 1
				}
				newChunk = chunk
			case Timer:
				if chunk.QtyVal != NoQty {
					chunk.Qty = fmt.Sprintf("%v", chunk.QtyVal)
				}
				newChunk = chunk
			default:
				t.Fatalf("Unable to process chunk: %v\n", chunk)
			}
			got.Steps[i][j] = newChunk
		}
	}

	// Now we can compare normally
	if !reflect.DeepEqual(*want, *got) {
		t.Fatalf("Assertion failed:\ngot:\t%+v\nwant:\t%+v", *got, *want)
	}
}

// --------------------------------------------------------------
// Tests
// --------------------------------------------------------------

func TestDirectionsWithDegrees(t *testing.T) {
	got := ParseRecipeString("", `Heat oven up to 200°C
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Heat oven up to 200°C")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestTimerWithName(t *testing.T) {
	got := ParseRecipeString("", `Fry for ~potato{42%minutes}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "potato", Qty: "42", QtyVal: 42, Unit: "minutes"}}, Steps: []Step{{Text("Fry for "), Timer{Name: "potato", Qty: "42", QtyVal: 42, Unit: "minutes"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidMultiWordTimer(t *testing.T) {
	got := ParseRecipeString("", `It is ~ {5}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("It is ~ {5}")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordIngredientWithPunctuation(t *testing.T) {
	got := ParseRecipeString("", `Add some @chilli, then serve
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add some "), Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, Text(", then serve")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientWithNumbers(t *testing.T) {
	got := ParseRecipeString("", `@tipo 00 flour{250%g}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "tipo 00 flour", Qty: "250", QtyVal: 250, Unit: "g"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "tipo 00 flour", Qty: "250", QtyVal: 250, Unit: "g"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMetadataMultiwordKeyWithSpaces(t *testing.T) {
	got := ParseRecipeString("", `>>cooking time    :30 mins
`)
	want := Recipe{Name: "", Metadata: map[string]string{"cooking time": "30 mins"}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidSingleWordIngredient(t *testing.T) {
	got := ParseRecipeString("", `Message me @ example
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Message me @ example")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentOneWord(t *testing.T) {
	got := ParseRecipeString("", `Simmer in #pan for some time
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "pan", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Simmer in "), Cookware{Name: "pan", Qty: "1", QtyVal: 1, Unit: ""}, Text(" for some time")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestFractionsLike(t *testing.T) {
	got := ParseRecipeString("", `@milk{01/2%cup}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "milk", Qty: "01/2", QtyVal: NoQty, Unit: "cup"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "milk", Qty: "01/2", QtyVal: NoQty, Unit: "cup"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientExplicitUnitsWithSpaces(t *testing.T) {
	got := ParseRecipeString("", `@chilli{ 3 % items }
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMultiWordIngredientNoAmount(t *testing.T) {
	got := ParseRecipeString("", `@hot chilli{}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "hot chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "hot chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordTimer(t *testing.T) {
	got := ParseRecipeString("", `Let it ~rest after plating
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}}, Steps: []Step{{Text("Let it "), Timer{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}, Text(" after plating")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentQuantityOneWord(t *testing.T) {
	got := ParseRecipeString("", `#frying pan{three}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "frying pan", Qty: "three", QtyVal: NoQty, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Cookware{Name: "frying pan", Qty: "three", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestFractions(t *testing.T) {
	got := ParseRecipeString("", `@milk{1/2%cup}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "milk", Qty: "0.5", QtyVal: 0.5, Unit: "cup"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "milk", Qty: "0.5", QtyVal: 0.5, Unit: "cup"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientMultipleWordsWithLeadingNumber(t *testing.T) {
	got := ParseRecipeString("", `Top with @1000 island dressing{ }
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "1000 island dressing", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Top with "), Ingredient{Name: "1000 island dressing", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMultiLineDirections(t *testing.T) {
	got := ParseRecipeString("", `Add a bit of chilli

Add a bit of hummus
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add a bit of chilli")}, {Text("Add a bit of hummus")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestTimerInteger(t *testing.T) {
	got := ParseRecipeString("", `Fry for ~{10%minutes}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "", Qty: "10", QtyVal: 10, Unit: "minutes"}}, Steps: []Step{{Text("Fry for "), Timer{Name: "", Qty: "10", QtyVal: 10, Unit: "minutes"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestComments(t *testing.T) {
	got := ParseRecipeString("", `-- testing comments
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestDirectionWithIngredient(t *testing.T) {
	got := ParseRecipeString("", `Add @chilli{3%items}, @ginger{10%g} and @milk{1%l}.
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}, {Name: "ginger", Qty: "10", QtyVal: 10, Unit: "g"}, {Name: "milk", Qty: "1", QtyVal: 1, Unit: "l"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add "), Ingredient{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}, Text(", "), Ingredient{Name: "ginger", Qty: "10", QtyVal: 10, Unit: "g"}, Text(" and "), Ingredient{Name: "milk", Qty: "1", QtyVal: 1, Unit: "l"}, Text(".")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMetadataMultiwordKey(t *testing.T) {
	got := ParseRecipeString("", `>> cooking time: 30 mins
`)
	want := Recipe{Name: "", Metadata: map[string]string{"cooking time": "30 mins"}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidSingleWordCookware(t *testing.T) {
	got := ParseRecipeString("", `Recipe # 5
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Recipe # 5")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestFractionsWithSpaces(t *testing.T) {
	got := ParseRecipeString("", `@milk{1 / 2 %cup}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "milk", Qty: "0.5", QtyVal: 0.5, Unit: "cup"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "milk", Qty: "0.5", QtyVal: 0.5, Unit: "cup"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMultiWordIngredient(t *testing.T) {
	got := ParseRecipeString("", `@hot chilli{3}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "hot chilli", Qty: "3", QtyVal: 3, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "hot chilli", Qty: "3", QtyVal: 3, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestQuantityDigitalString(t *testing.T) {
	got := ParseRecipeString("", `@water{7 k }
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "water", Qty: "7 k", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "water", Qty: "7 k", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordIngredientWithUnicodePunctuation(t *testing.T) {
	got := ParseRecipeString("", `Add @chilli⸫ then bake
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add "), Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, Text("⸫ then bake")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidMultiWordCookware(t *testing.T) {
	got := ParseRecipeString("", `Recipe # 10{}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Recipe # 10{}")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestBasicDirection(t *testing.T) {
	got := ParseRecipeString("", `Add a bit of chilli
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add a bit of chilli")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentMultipleWords(t *testing.T) {
	got := ParseRecipeString("", `Fry in #frying pan{}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "frying pan", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Fry in "), Cookware{Name: "frying pan", Qty: "1", QtyVal: 1, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMetadataBreak(t *testing.T) {
	got := ParseRecipeString("", `hello >> sourced: babooshka
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("hello >> sourced: babooshka")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMultipleLines(t *testing.T) {
	got := ParseRecipeString("", `>> Prep Time: 15 minutes
>> Cook Time: 30 minutes
`)
	want := Recipe{Name: "", Metadata: map[string]string{"Prep Time": "15 minutes", "Cook Time": "30 minutes"}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordCookwareWithUnicodePunctuation(t *testing.T) {
	got := ParseRecipeString("", `Place in #pot⸫ then boil
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Place in "), Cookware{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}, Text("⸫ then boil")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestCookwareWithUnicodeWhitespace(t *testing.T) {
	got := ParseRecipeString("", `Add to #pot then boil
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Add to "), Cookware{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}, Text(" then boil")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSlashInText(t *testing.T) {
	got := ParseRecipeString("", `Preheat the oven to 200℃/Fan 180°C.
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Preheat the oven to 200℃/Fan 180°C.")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordTimerWithUnicodePunctuation(t *testing.T) {
	got := ParseRecipeString("", `Let it ~rest⸫ then serve
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}}, Steps: []Step{{Text("Let it "), Timer{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}, Text("⸫ then serve")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentMultipleWordsWithLeadingNumber(t *testing.T) {
	got := ParseRecipeString("", `Fry in #7-inch nonstick frying pan{ }
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "7-inch nonstick frying pan", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Fry in "), Cookware{Name: "7-inch nonstick frying pan", Qty: "1", QtyVal: 1, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentQuantityMultipleWords(t *testing.T) {
	got := ParseRecipeString("", `#frying pan{two small}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "frying pan", Qty: "two small", QtyVal: NoQty, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Cookware{Name: "frying pan", Qty: "two small", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordTimerWithPunctuation(t *testing.T) {
	got := ParseRecipeString("", `Let it ~rest, then serve
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}}, Steps: []Step{{Text("Let it "), Timer{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}, Text(", then serve")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestTimerWithUnicodeWhitespace(t *testing.T) {
	got := ParseRecipeString("", `Let it ~rest then serve
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}}, Steps: []Step{{Text("Let it "), Timer{Name: "rest", Qty: "", QtyVal: NoQty, Unit: ""}, Text(" then serve")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentMultipleWordsWithSpaces(t *testing.T) {
	got := ParseRecipeString("", `Fry in #frying pan{ }
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "frying pan", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Fry in "), Cookware{Name: "frying pan", Qty: "1", QtyVal: 1, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientExplicitUnits(t *testing.T) {
	got := ParseRecipeString("", `@chilli{3%items}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "3", QtyVal: 3, Unit: "items"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientNoUnitsNotOnlyString(t *testing.T) {
	got := ParseRecipeString("", `@5peppers
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "5peppers", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "5peppers", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestServings(t *testing.T) {
	got := ParseRecipeString("", `>> servings: 1|2|3
`)
	want := Recipe{Name: "", Metadata: map[string]string{"servings": "1|2|3"}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidMultiWordIngredient(t *testing.T) {
	got := ParseRecipeString("", `Message @ example{}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Message @ example{}")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestCommentsWithIngredients(t *testing.T) {
	got := ParseRecipeString("", `-- testing comments
@thyme{2%sprigs}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "thyme", Qty: "2", QtyVal: 2, Unit: "sprigs"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "thyme", Qty: "2", QtyVal: 2, Unit: "sprigs"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientNoUnits(t *testing.T) {
	got := ParseRecipeString("", `@chilli
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestQuantityAsText(t *testing.T) {
	got := ParseRecipeString("", `@thyme{few%sprigs}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "thyme", Qty: "few", QtyVal: NoQty, Unit: "sprigs"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "thyme", Qty: "few", QtyVal: NoQty, Unit: "sprigs"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientImplicitUnits(t *testing.T) {
	got := ParseRecipeString("", `@chilli{3}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "3", QtyVal: 3, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "3", QtyVal: 3, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestTimerFractional(t *testing.T) {
	got := ParseRecipeString("", `Fry for ~{1/2%hour}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "", Qty: "0.5", QtyVal: 0.5, Unit: "hour"}}, Steps: []Step{{Text("Fry for "), Timer{Name: "", Qty: "0.5", QtyVal: 0.5, Unit: "hour"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestInvalidSingleWordTimer(t *testing.T) {
	got := ParseRecipeString("", `It is ~ 5
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("It is ~ 5")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientWithUnicodeWhitespace(t *testing.T) {
	got := ParseRecipeString("", `Add @chilli then bake
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add "), Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, Text(" then bake")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestSingleWordCookwareWithPunctuation(t *testing.T) {
	got := ParseRecipeString("", `Place in #pot, then boil
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Text("Place in "), Cookware{Name: "pot", Qty: "1", QtyVal: 1, Unit: ""}, Text(", then boil")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestDirectionsWithNumbers(t *testing.T) {
	got := ParseRecipeString("", `Heat 5L of water
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Heat 5L of water")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMetadata(t *testing.T) {
	got := ParseRecipeString("", `>> sourced: babooshka
`)
	want := Recipe{Name: "", Metadata: map[string]string{"sourced": "babooshka"}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestEquipmentQuantity(t *testing.T) {
	got := ParseRecipeString("", `#frying pan{2}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{{Name: "frying pan", Qty: "2", QtyVal: 2, Unit: ""}}, Timers: []Timer{}, Steps: []Step{{Cookware{Name: "frying pan", Qty: "2", QtyVal: 2, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestFractionsInDirections(t *testing.T) {
	got := ParseRecipeString("", `knife cut about every 1/2 inches
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("knife cut about every 1/2 inches")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientWithEmoji(t *testing.T) {
	got := ParseRecipeString("", `Add some @🧂
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "🧂", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Text("Add some "), Ingredient{Name: "🧂", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestIngredientWithoutStopper(t *testing.T) {
	got := ParseRecipeString("", `@chilli cut into pieces
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, Text(" cut into pieces")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestMutipleIngredientsWithoutStopper(t *testing.T) {
	got := ParseRecipeString("", `@chilli cut into pieces and @garlic
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, {Name: "garlic", Qty: "some", QtyVal: NoQty, Unit: ""}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "chilli", Qty: "some", QtyVal: NoQty, Unit: ""}, Text(" cut into pieces and "), Ingredient{Name: "garlic", Qty: "some", QtyVal: NoQty, Unit: ""}}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestCommentsAfterIngredients(t *testing.T) {
	got := ParseRecipeString("", `@thyme{2%sprigs} -- testing comments
and some text
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{{Name: "thyme", Qty: "2", QtyVal: 2, Unit: "sprigs"}}, Cookware: []Cookware{}, Timers: []Timer{}, Steps: []Step{{Ingredient{Name: "thyme", Qty: "2", QtyVal: 2, Unit: "sprigs"}, Text(" ")}, {Text("and some text")}}}
	assertCanonicalRecipe(t, &got, &want)
}
func TestTimerDecimal(t *testing.T) {
	got := ParseRecipeString("", `Fry for ~{1.5%minutes}
`)
	want := Recipe{Name: "", Metadata: map[string]string{}, Ingredients: []Ingredient{}, Cookware: []Cookware{}, Timers: []Timer{{Name: "", Qty: "1.5", QtyVal: 1.5, Unit: "minutes"}}, Steps: []Step{{Text("Fry for "), Timer{Name: "", Qty: "1.5", QtyVal: 1.5, Unit: "minutes"}}}}
	assertCanonicalRecipe(t, &got, &want)
}
