package trmnl

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	promUrl = os.Getenv("PROMETHEUS_URI")
)

func FetchLast48Hours(query string) []float64 {
	end := time.Now().Unix()
	start := end - (2 * 24 * 3600)
	resp, err := http.Get(promUrl + "/api/v1/query_range" +
		"?query=" + query +
		"&step=300" +
		"&start=" + strconv.FormatInt(start, 10) +
		"&end=" + strconv.FormatInt(end, 10))
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var x map[string]interface{}
	if err := json.Unmarshal(body, &x); err != nil {
		log.Fatal(err)
	}
	if (x["status"]) != "success" {
		log.Fatal("Prom response not successful")
	}
	data := x["data"].(map[string]interface{})
	result := data["result"].([]interface{})
	if len(result) != 1 {
		log.Fatal("Expected 1 result, got " + strconv.Itoa(len(result)))
	}
	firstResult := result[0].(map[string]interface{})
	values := firstResult["values"].([]interface{})
	var temperatures = make([]float64, len(values))
	for i := 0; i < len(values); i++ {
		value := values[i].([]interface{})
		temperatures[i], err = strconv.ParseFloat(value[1].(string), 64)
	}
	return temperatures
}
