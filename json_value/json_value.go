package json_value

func ToInt(jsonValue interface{}, path string,
	usagef func(string, ...interface{})) int {

	number, ok := jsonValue.(float64)
	if !ok {
		usagef("%s must be int", path)
	}
	if float64(int(number)) != number {
		usagef("%s must be int", path)
	}
	return int(number)
}

func ToString(jsonValue interface{}, path string, usagef func(string, ...interface{})) string {
	str, ok := jsonValue.(string)
	if !ok {
		usagef("%s must be string", path)
	}
	return str
}

func ToMap(jsonValue interface{}, path string, usagef func(string, ...interface{})) map[string]interface{} {
	m, ok := jsonValue.(map[string]interface{})
	if !ok {
		usagef("%s must be object", path)
	}
	return m
}
