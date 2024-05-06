package train

import (
	"encoding/json"
	"os"
)

func readJSON(fileName string) map[string]int {
	datas := map[string]int{}

	file, _ := os.ReadFile(fileName)
	json.Unmarshal(file, &datas)

	return datas
}

func writeJSON(fileName string, data map[string]int) error {
	jsonString, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, jsonString, os.ModePerm)
}
