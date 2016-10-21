# goa-json2type

goa-json2type convert JSON format string to goa Type like below.

sample.json

```json
{
    "user": {
        "name": "hatajoe",
        "email": "foo@example.com",
        "address": {
            "country": "Japan",
            "city": "Osaka"
        },
        "languages": ["japanese", "go", "js"],
        "age": 33,
        "admin": true
    },
    "tag": "engineers"
}
```

```sh
% cat sample.json | goa-json2type -t media
// 1
var address = Type("address", func() {
	
	Member("country", String, "", func() {})
	Member("city", String, "", func() {})
})
// 2
var user = Type("user", func() {
	
	Member("admin", Boolean, "", func() {})
	Member("name", String, "", func() {})
	Member("email", String, "", func() {})
	Member("address", address, "", func() {})
	Member("languages", ArrayOf(String), "", func() {})
	Member("age", Integer, "", func() {})
})

// 3
func() {

	Attribute("user", user, "", func() {})
	Attribute("tag", String, "", func() {})
    // if `-t payload' specified ...
	// Member("user", user, "", func() {})
	// Member("tag", String, "", func() {})

	Required("user")
	Required("tag")
}
```

```go
// design/types.go

// paste 1, 2 here
```

```go
// design/media_types.go

package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var SomeMedia = MediaType("application/vnd.example.some+json", func() {
	Description("blah blah blah...")
	Attributes(
        // paste 3 here 
    )
	View("default", func() {
		Attribute("user")
		Attribute("tag")
	})
})
```

## options 

```
  -t string
    	media or payload (default "media")
```

# LICENCE

MIT
