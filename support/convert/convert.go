package convert

import "encoding/json"

func Struct2Json2Map(obj interface{}) (map[string]interface{}, error) {
	var converted = make(map[string]interface{})
	jsonByte, err := json.Marshal(obj)
	if err != nil {
		return converted, err
	}
	if err = json.Unmarshal(jsonByte, &converted); err != nil {
		return converted, err
	}
	return converted, nil
}
