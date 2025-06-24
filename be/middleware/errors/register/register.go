package register

var (
	defaultLang = "default"
	codes       = make(map[string]map[int32]string, 0)
)

// Register 不支持并发处理
func Register(lang string, partCodes map[int32]string) {
	if engCodes := codes[lang]; engCodes == nil {
		codes[lang] = make(map[int32]string, 0)
	}
	for code, message := range partCodes {
		if _, ok := codes[lang][code]; ok {
			panic("error code duplicate. code: %d, message: %s")
		}
		codes[lang][code] = message
	}

}

func Get(lang string, code int32) string {
	if langCodes, ok := codes[lang]; ok {
		if message, ok := langCodes[code]; ok {
			return message
		}
	}

	if langCodes, ok := codes[defaultLang]; ok {
		if message, ok := langCodes[code]; ok {
			return message
		}
	}

	return ""
}
