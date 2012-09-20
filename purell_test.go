package purell

import (
	"github.com/rogpeppe/purell"
	"testing"
)

var tests = []struct {
	url    string
	flags  purell.NormalizationFlags
	expect string
}{{
	"HTTP://www.SRC.ca",
	purell.FlagLowercaseScheme,
	"http://www.SRC.ca",
}, {
	"HTTP://www.SRC.ca",
	purell.FlagLowercaseScheme,
	"http://www.SRC.ca",
}, {
	"http://www.SRC.ca",
	purell.FlagLowercaseScheme,
	"http://www.SRC.ca",
}, {
	"HTTP://www.SRC.ca/",
	purell.FlagLowercaseHost,
	"HTTP://www.src.ca/",
}, {
	"http://www.whatever.com/Some%aa%20Special%8Ecases/",
	purell.FlagUppercaseEscapes,
	"http://www.whatever.com/Some%AA%20Special%8Ecases/",
}, {
	"http://www.toto.com/%41%42%2E%44/%32%33%52%2D/%5f%7E",
	purell.FlagDecodeUnnecessaryEscapes,
	"http://www.toto.com/AB.D/23R-/_~",
}, {
	"HTTP://www.SRC.ca:80/",
	purell.FlagRemoveDefaultPort,
	"HTTP://www.SRC.ca/",
}, {
	"HTTP://www.SRC.ca:80",
	purell.FlagRemoveDefaultPort,
	"HTTP://www.SRC.ca",
}, {
	"HTTP://www.SRC.ca:8080",
	purell.FlagRemoveDefaultPort,
	"HTTP://www.SRC.ca:8080",
}, {
	"HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e",
	purell.FlagsSafe,
	"http://www.src.ca/to%1Ato%8B%EE/OKnowABC~",
}, {
	"HTTP://www.SRC.ca:80/to%1ato%8b%ee/OKnow%41%42%43%7e",
	purell.FlagLowercaseHost | purell.FlagLowercaseScheme,
	"http://www.src.ca:80/to%1Ato%8B%EE/OKnowABC~",
}, {
	"HTTP://www.SRC.ca:80/",
	purell.FlagRemoveTrailingSlash,
	"HTTP://www.SRC.ca:80",
}, {
	"HTTP://www.SRC.ca:80/toto/titi/",
	purell.FlagRemoveTrailingSlash,
	"HTTP://www.SRC.ca:80/toto/titi",
}, {
	"HTTP://www.SRC.ca:80/toto/titi/fin/?a=1",
	purell.FlagRemoveTrailingSlash,
	"HTTP://www.SRC.ca:80/toto/titi/fin?a=1",
}, {
	"HTTP://www.SRC.ca:80",
	purell.FlagAddTrailingSlash,
	"HTTP://www.SRC.ca:80/",
}, {
	"HTTP://www.SRC.ca:80/toto/titi.html",
	purell.FlagAddTrailingSlash,
	"HTTP://www.SRC.ca:80/toto/titi.html/",
}, {
	"HTTP://www.SRC.ca:80/toto/titi/fin?a=1",
	purell.FlagAddTrailingSlash,
	"HTTP://www.SRC.ca:80/toto/titi/fin/?a=1",
}, {
	"HTTP://root/a/b/./../../c/",
	purell.FlagRemoveDotSegments,
	"HTTP://root/c/",
}, {
	"HTTP://root/../a/b/./../c/../d",
	purell.FlagRemoveDotSegments,
	"HTTP://root/a/d",
}, {
	"HTTP://www.SRC.ca:80/to%1ato%8b%ee/./c/d/../OKnow%41%42%43%7e/?a=b#test",
	purell.FlagsUsuallySafe,
	"http://www.src.ca/to%1Ato%8B%EE/c/OKnowABC~?a=b#test",
}, {
	"HTTP://root/a/b/c/default.aspx",
	purell.FlagRemoveDirectoryIndex,
	"HTTP://root/a/b/c/",
}, {
	"HTTP://root/a/b/c/default#a=b",
	purell.FlagRemoveDirectoryIndex,
	"HTTP://root/a/b/c/default#a=b",
}, {
	"HTTP://root/a/b/c/default#toto=tata",
	purell.FlagRemoveFragment,
	"HTTP://root/a/b/c/default",
}, {
	"https://root/a/b/c/default#toto=tata",
	purell.FlagForceHttp,
	"http://root/a/b/c/default#toto=tata",
}, {
	"https://root/a//b///c////default#toto=tata",
	purell.FlagRemoveDuplicateSlashes,
	"https://root/a/b/c/default#toto=tata",
}, {
	"https://root//a//b///c////default#toto=tata",
	purell.FlagRemoveDuplicateSlashes,
	"https://root/a/b/c/default#toto=tata",
}, {
	"https://www.root/a/b/c/",
	purell.FlagRemoveWWW,
	"https://root/a/b/c/",
}, {
	"https://WwW.Root/a/b/c/",
	purell.FlagRemoveWWW,
	"https://Root/a/b/c/",
}, {
	"https://Root/a/b/c/",
	purell.FlagAddWWW,
	"https://www.Root/a/b/c/",
}, {
	"http://root/toto/?b=4&a=1&c=3&b=2&a=5",
	purell.FlagSortQuery,
	"http://root/toto/?a=1&a=5&b=2&b=4&c=3",
}, {
	"http://root/toto/?",
	purell.FlagRemoveEmptyQuerySeparator,
	"http://root/toto/",
}, {
	"HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid",
	purell.FlagsUnsafe,
	"http://root.com/toto/tE%1F/a/c?a=4&w=1&w=2&z=3",
}, {
	"HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid",
	purell.FlagsSafe,
	"https://www.root.com/toto/tE%1F///a/./b/../c/?z=3&w=2&a=4&w=1#invalid",
}, {
	"HTTPS://www.RooT.com/toto/t%45%1f///a/./b/../c/?z=3&w=2&a=4&w=1#invalid",
	purell.FlagsUsuallySafe,
	"https://www.root.com/toto/tE%1F///a/c?z=3&w=2&a=4&w=1#invalid",
},
}

func TestNormalize(t *testing.T) {
	for _, test := range tests {
		got, err := purell.NormalizeURLString(test.url, test.flags)
		if err != nil {
			t.Errorf("got error on %q: %v", test.url, err)
		} else if got != test.expect {
			t.Errorf("normalizing url %q, flags %v: expected %q; got %q", test.url, test.flags, test.expect, got)
		}
	}
}
