package context

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/rentiansheng/dashboard/middleware/context/decode"
	"github.com/rentiansheng/dashboard/middleware/errors"
)

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEMSGPACK           = "application/x-msgpack"
	MIMEMSGPACK2          = "application/msgpack"
	MIMEYAML              = "application/x-yaml"
	MIMETOML              = "application/toml"
)

var (
	validate = validator.New()
)

type ValidateI interface {
	Validate() errors.Error
}

type ValidateRawI interface {
	Validate() error
}

type DefaultI interface {
	Default()
}

// autoDecode  自动适配需要解析方式。已经支持自动解析query string, post(application/x-www-form-urlencoded,json)
// 默认是先解析 query string ,然后根据 http content type 解析 合并数据。默认是json 方式
func autoDecode(req *http.Request, urlParams map[string][]string, obj interface{}) error {

	if err := decode.Query(req, obj); err != nil {
		return err
	}

	contextTypes := strings.Split(req.Header.Get("Content-Type"), ";")
	switch contextTypes[0] {
	case MIMEPOSTForm:
		if err := decode.Form(req, obj); err != nil {
			return err
		}
	case MIMEJSON:
		reqBodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		if len(reqBodyBytes) > 0 {
			decoder := json.NewDecoder(bytes.NewReader(reqBodyBytes))
			if EnableDecoderUseNumber {
				decoder.UseNumber()
			}
			if EnableDecoderDisallowUnknownFields {
				decoder.DisallowUnknownFields()
			}

			if err := decoder.Decode(obj); err != nil {
				return err
			}
		}

	}
	newHeader := make(map[string][]string, 0)
	for c, v := range req.Header {
		newHeader[strings.ToLower(c)] = v
	}
	if err := decode.HTTPHeader(newHeader, obj); err != nil {
		return err
	}
	if err := decode.HTTPUri(urlParams, obj); err != nil {
		return err
	}

	if defaultSet, ok := obj.(DefaultI); ok {
		defaultSet.Default()
	}

	if validate == nil {
		return nil
	}
	if err := validateStruct(obj); err != nil {
		return err
	}

	if valid, ok := obj.(ValidateRawI); ok {
		if err := valid.Validate(); err != nil {
			return err
		}
	}

	return nil
}

/*
Fields:
Tag	Description
eqcsfield	Field Equals Another Field (relative)
eqfield	Field Equals Another Field
fieldcontains	NOT DOCUMENTED IN doc.go
fieldexcludes	NOT DOCUMENTED IN doc.go
gtcsfield	Field Greater Than Another Relative Field
gtecsfield	Field Greater Than or Equal To Another Relative Field
gtefield	Field Greater Than or Equal To Another Field
gtfield	Field Greater Than Another Field
ltcsfield	Less Than Another Relative Field
ltecsfield	Less Than or Equal To Another Relative Field
ltefield	Less Than or Equal To Another Field
ltfield	Less Than Another Field
necsfield	Field Does Not Equal Another Field (relative)
nefield	Field Does Not Equal Another Field
Network:
Tag	Description
cidr	Classless Inter-Domain Routing CIDR
cidrv4	Classless Inter-Domain Routing CIDRv4
cidrv6	Classless Inter-Domain Routing CIDRv6
datauri	Data URL
fqdn	Full Qualified Domain Name (FQDN)
hostname	Hostname RFC 952
hostname_port	HostPort
hostname_rfc1123	Hostname RFC 1123
ip	Internet Protocol Address IP
ip4_addr	Internet Protocol Address IPv4
ip6_addr	Internet Protocol Address IPv6
ip_addr	Internet Protocol Address IP
ipv4	Internet Protocol Address IPv4
ipv6	Internet Protocol Address IPv6
mac	Media Access Control Address MAC
tcp4_addr	Transmission Control Protocol Address TCPv4
tcp6_addr	Transmission Control Protocol Address TCPv6
tcp_addr	Transmission Control Protocol Address TCP
udp4_addr	User Datagram Protocol Address UDPv4
udp6_addr	User Datagram Protocol Address UDPv6
udp_addr	User Datagram Protocol Address UDP
unix_addr	Unix domain socket end point Address
uri	URI String
url	URL String
url_encoded	URL Encoded
urn_rfc2141	Urn RFC 2141 String
Strings:
Tag	Description
alpha	Alpha Only
alphanum	Alphanumeric
alphanumunicode	Alphanumeric Unicode
alphaunicode	Alpha Unicode
ascii	ASCII
boolean	Boolean
contains	Contains
containsany	Contains Any
containsrune	Contains Rune
endsnotwith	Ends With
endswith	Ends With
excludes	Excludes
excludesall	Excludes All
excludesrune	Excludes Rune
lowercase	Lowercase
multibyte	Multi-Byte Characters
number	NOT DOCUMENTED IN doc.go
numeric	Numeric
printascii	Printable ASCII
startsnotwith	Starts Not With
startswith	Starts With
uppercase	Uppercase
Format:
Tag	Description
base64	Base64 String
base64url	Base64URL String
bic	Business Identifier Code (ISO 9362)
bcp47_language_tag	Language tag (BCP 47)
btc_addr	Bitcoin Address
btc_addr_bech32	Bitcoin Bech32 Address (segwit)
datetime	Datetime
e164	e164 formatted phone number
email	E-mail String
eth_addr	Ethereum Address
hexadecimal	Hexadecimal String
hexcolor	Hexcolor String
hsl	HSL String
hsla	HSLA String
html	HTML Tags
html_encoded	HTML Encoded
isbn	International Standard Book Number
isbn10	International Standard Book Number 10
isbn13	International Standard Book Number 13
iso3166_1_alpha2	Two-letter country code (ISO 3166-1 alpha-2)
iso3166_1_alpha3	Three-letter country code (ISO 3166-1 alpha-3)
iso3166_1_alpha_numeric	Numeric country code (ISO 3166-1 numeric)
iso3166_2	Country subdivision code (ISO 3166-2)
iso4217	Currency code (ISO 4217)
json	JSON
jwt	JSON Web Token (JWT)
latitude	Latitude
longitude	Longitude
postcode_iso3166_alpha2	Postcode
postcode_iso3166_alpha2_field	Postcode
rgb	RGB String
rgba	RGBA String
ssn	Social Security Number SSN
timezone	Timezone
uuid	Universally Unique Identifier UUID
uuid3	Universally Unique Identifier UUID v3
uuid3_rfc4122	Universally Unique Identifier UUID v3 RFC4122
uuid4	Universally Unique Identifier UUID v4
uuid4_rfc4122	Universally Unique Identifier UUID v4 RFC4122
uuid5	Universally Unique Identifier UUID v5
uuid5_rfc4122	Universally Unique Identifier UUID v5 RFC4122
uuid_rfc4122	Universally Unique Identifier UUID RFC4122
semver	Semantic Versioning 2.0.0
ulid	Universally Unique Lexicographically Sortable Identifier ULID
Comparisons:
Tag	Description
eq	Equals
gt	Greater than
gte	Greater than or equal
lt	Less Than
lte	Less Than or Equal
ne	Not Equal
Other:
Tag	Description
dir	Directory
file	File path
isdefault	Is Default
len	Length
max	Maximum
min	Minimum
oneof	One Of
required	Required
required_if	Required If
required_unless	Required Unless
required_with	Required With
required_with_all	Required With All
required_without	Required Without
required_without_all	Required Without All
excluded_with	Excluded With
excluded_with_all	Excluded With All
excluded_without	Excluded Without
excluded_without_all	Excluded Without All
unique	Unique
Aliases:
Tag	Description
iscolor	hexcolor|rgb|rgba|hsl|hsla
country_code	iso3166_1_alpha2|iso3166_1_alpha3|iso3166_1_alpha_numeric
*/
func decodeJSON(r io.Reader, obj interface{}) error {
	reqBodyBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if len(reqBodyBytes) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(reqBodyBytes))
		if EnableDecoderUseNumber {
			decoder.UseNumber()
		}
		if EnableDecoderDisallowUnknownFields {
			decoder.DisallowUnknownFields()
		}

		if err := decoder.Decode(obj); err != nil {
			return err
		}

	}

	if defaultSet, ok := obj.(DefaultI); ok {
		defaultSet.Default()
	}

	if validate == nil {
		return nil
	}
	if err := validateStruct(obj); err != nil {
		return err
	}

	if valid, ok := obj.(ValidateRawI); ok {
		if err := valid.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func validateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return validate.Struct(value.Elem().Interface())
	case reflect.Struct:
		return validate.Struct(obj)
	case reflect.Slice, reflect.Array:
		return fmt.Errorf("slice unimplement")
	default:
		return nil
	}

}

/*
func decodeForm(req *http.Request, obj interface{}) error {
	if err := decode.Form(req, obj); err != nil {
		return err
	}
	if validate == nil {
		return nil
	}
	if err := validateStruct(obj); err != nil {
		return err
	}

	if valid, ok := obj.(ValidateRawI); ok {
		if err := valid.Validate(); err != nil {
			return err
		}
	}

	return nil
}
*/
