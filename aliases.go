package stubs

var (
	generatorAliases map[string]string
)

func init() {
	// registers aliases for generator names
	RegisterAltGenNames("bool", "boolean", "flag")
	RegisterAltGenNames("company", "company-name", "corporation", "business")
	RegisterAltGenNames("company-bs", "company-mission")
	RegisterAltGenNames("company-slogan", "company-catch-phrase")
	RegisterAltGenNames("country", "country-name")
	RegisterAltGenNames("credit-card", "creditcard")
	RegisterAltGenNames("domain", "domain-name")
	RegisterAltGenNames("hexcolor", "hex-color", "hexcolour", "hex-colour")
	RegisterAltGenNames("hostname", "domainword", "domain-word", "host", "host-name")
	RegisterAltGenNames("ipv4", "ip4", "ip", "ip-address")
	RegisterAltGenNames("ipv6", "ip6")
	RegisterAltGenNames("isbn10", "isbnv10", "isbn")
	RegisterAltGenNames("isbn13", "isbnv13")
	RegisterAltGenNames("landline", "phone-number", "phone", "telephone")
	RegisterAltGenNames("latitude", "lat")
	RegisterAltGenNames("longitude", "lon")
	RegisterAltGenNames("mac-address", "mac", "macaddress")
	RegisterAltGenNames("mobile", "mobile-number", "cell", "cell-phone", "gsm", "gsm-number")
	RegisterAltGenNames("postcode", "zipcode", "post-code", "zip-code", "zip")
	RegisterAltGenNames("rgbcolor", "rgb-color", "rgbcolour", "rgb-colour")
	RegisterAltGenNames("sentences", "text", "phrases")
	RegisterAltGenNames("ssn", "socialsecurity", "social-security", "social-security-number", "ss-number")
	RegisterAltGenNames("state", "state-code")
	RegisterAltGenNames("user-name", "username", "login", "nickname", "nick-name")
	RegisterAltGenNames("uuid3", "uuidv3")
	RegisterAltGenNames("uuid4", "uuidv4")
	RegisterAltGenNames("uuid5", "uuidv5")
	RegisterAltGenNames("float", "float32")
	RegisterAltGenNames("double", "float64")
	RegisterAltGenNames("datetime", "date-time")
	RegisterAltGenNames("amount", "price", "currency-amount", "cost", "turnover", "vat")
	RegisterAltGenNames("small-amount", "small-price", "low-price", "small-currency-amount", "low-cost", "fees")

	// Now add keys for all known generators
	g, _ := newGenerator("")
	for key := range g.gens {
		if generatorAliases[key] == "" {
			RegisterAltGenNames(key, key)
		}
	}
}

// RegisterAltGenNames registers alternatives for a generator name
func RegisterAltGenNames(key string, alts ...string) {
	if generatorAliases == nil {
		generatorAliases = make(map[string]string, 300)
	}
	for _, v := range alts {
		generatorAliases[v] = key
	}
}
