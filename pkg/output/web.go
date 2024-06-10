package output

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// var trackFile = "scoring/scored.csv"

type ScoreRequest struct {
	Domain string `json:"domain"`
	Score  int    `json:"score"`
}

func (engine *RenderEngine) RunHttp(dir string) {

	fmt.Println("Starting output http server...")

	http.HandleFunc("/score", engine.handleScore)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	http.ListenAndServe(":8081", nil)
}

func (engine *RenderEngine) handleScore(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := io.ReadAll(r.Body)

	var data ScoreRequest

	err := json.Unmarshal(reqBody, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	existingCandidates, err := ioutil.ReadFile("output/candidates.txt")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	candidates := strings.Split(string(existingCandidates), "\n")

	found := false
	for _, c := range candidates {
		if c == fmt.Sprintf("https://%s", data.Domain) {
			found = true
		}
	}

	if !found {
		f, err := os.OpenFile("output/candidates.txt",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(fmt.Sprintf("https://%s\n", data.Domain)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	labels := readJSON("output/labels.json")

	labels[data.Domain] = data.Score
	writeJSON("output/labels.json", labels)

}

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
