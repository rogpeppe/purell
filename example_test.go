package purell

import (
	"fmt"
	"net/url"
)

func ExampleNormalizeURLString() {
	if normalized, err := NormalizeURLString("hTTp://someWEBsite.com:80/Amazing%3f/url/",
		FlagLowercaseScheme|FlagLowercaseHost|FlagUppercaseEscapes); err != nil {
		panic(err)
	} else {
		fmt.Print(normalized)
	}
	// Output: http://somewebsite.com:80/Amazing%3F/url/
}

func ExampleMustNormalizeURLString() {
	normalized := MustNormalizeURLString("hTTpS://someWEBsite.com:80/Amazing%fa/url/",
		FlagsUnsafe)
	fmt.Print(normalized)

	// Output: http://somewebsite.com/Amazing%FA/url
}

func ExampleNormalizeURL() {
	u, err := url.Parse("Http://SomeURL.com:8080/a/b/.././c///g?c=3&a=1&b=9&c=0#target")
	if err != nil {
		panic(err)
	}
	NormalizeURL(u, FlagsUsuallySafe|FlagRemoveDuplicateSlashes|FlagRemoveFragment)
	fmt.Print(u)

	// Output: http://someurl.com:8080/a/c/g?c=3&a=1&b=9&c=0
}
