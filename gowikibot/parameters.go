package gowikibot

type Values map[string]string

// Return the value associated with a specific key.
// if one was not found then we return an empty string.
func (val Values) Get(key string) string {

	if val == nil {
		return ""
	}

	valString, ok := val[key]

	if !ok {
		return ""
	}

	return valString
}

// set a value on a key and replace any exisitng value that may
// exist
func (val Values) Set(key, value string) {
	val[key] = value
}

// Delete the value associated with a key, obviously
func (val Values) Delete(key string) {
	delete(val, key)
}
