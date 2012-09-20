package purell

import (
	"fmt"
	"github.com/rogpeppe/purell"
	"net/url"
)

func ExampleNormalizeURLString() {
	if normalized, err := purell.NormalizeURLString("hTTp://someWEBsite.com:80/Amazing%3f/url/",
		purell.LowercaseScheme|purell.LowercaseHost|purell.UppercaseEscapes); err != nil {
		panic(err)
	} else {
		fmt.Print(normalized)
	}
	// Output: http://somewebsite.com:80/Amazing%3F/url/
}

func ExampleMustNormalizeURLString() {
	normalized := purell.MustNormalizeURLString("hTTpS://someWEBsite.com:80/Amazing%fa/url/",
		purell.sUnsafe)
	fmt.Print(normalized)

	// Output: http://somewebsite.com/Amazing%FA/url
}

func ExampleNormalizeURL() {
	u, err := url.Parse("Http://SomeURL.com:8080/a/b/.././c///g?c=3&a=1&b=9&c=0#target")
	if err != nil {
		panic(err)
	}
	NormalizeURL(u, purell.UsuallySafe|purell.RemoveDuplicateSlashes|purell.RemoveFragment)
	fmt.Print(u)

	// Output: http://someurl.com:8080/a/c/g?c=3&a=1&b=9&c=0
}
