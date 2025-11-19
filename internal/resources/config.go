package resources

// [[resources.http]]
// url = 'https://example.com'
// name = 'example'
type HttpResourceConfig struct {
	Url            string `toml:"url"`
	Name           string `toml:"name"`
	ExpectedStatus int    `toml:"expected_status"`
}

type ResourcesConfig struct {
	Http []HttpResourceConfig `toml:"http"`
}
