package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// const host = "http://localhost:5102/"
const host = "https://api.innobyte.dev/"

func main() {

	app := fiber.New()

	err := godotenv.Load()
	if err != nil {
		return
	}

	app.Get("/:key/:projectid/:endpointURL", func(c *fiber.Ctx) error {

		idProject := c.Params("projectid")
		key := c.Params("key")
		endpointURL := c.Params("endpointURL")

		logObject := map[string]interface{}{
			"type":            c.Method(),
			"status":          "200",
			"IdUser":          "",
			"IdProject":       idProject,
			"query":           "",
			"requestDuration": 1,
		}

		m := c.Queries()
		par := formatParam(m)
		jsonResponse, err, _ := apiRequest(host+"api/Project/details/"+idProject, c.Method(), nil)

		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logObject["IdUser"] = jsonResponse["idUser"]

		if key != jsonResponse["key"] {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("INCORRECT KEY PROOJECT")
		}
		endpointsRaw, ok := jsonResponse["endpoints"].([]interface{})
		if !ok {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("CANNOT DESERIALIZE ENDPOINTS")
		}
		//var endpoints []endpoint.Endpoint
		finded := false
		query := ""
		for _, e := range endpointsRaw {
			endpointMap, ok := e.(map[string]interface{})
			if !ok {
				logObject["status"] = "400"
				logRequest(logObject)
				return c.SendString("Invalid endpoint format")
			}

			if url, ok := endpointMap["url"].(string); ok {
				if url == endpointURL {
					finded = true
					query, ok = endpointMap["query"].(string)
					if !ok {
						query = ""
					}
				}
			}
		}

		if !finded {
			logObject["status"] = "404"
			logRequest(logObject)
			return c.SendString("COULD NOT FIND ANY ENDPOINT")
		}
		logObject["query"] = query

		data := map[string]interface{}{
			"database": jsonResponse["iddatabaseNavigation"],
			"query":    query,
			"params":   par,
		}

		dataRequest, err := json.Marshal(data)
		if err != nil {
			logObject["status"] = "500"
			logRequest(logObject)
			return c.SendString("ERROR")
		}

		_, err, jsonRes := apiRequest(host+"api/Request", "POST", dataRequest)
		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logRequest(logObject)
		return c.JSON(jsonRes)

	})

	app.Post("/:key/:projectid/:endpointURL", func(c *fiber.Ctx) error {
		idProject := c.Params("projectid")
		key := c.Params("key")
		endpointURL := c.Params("endpointURL")

		body := c.Body()
		var requestBody map[string]string

		if err := json.Unmarshal(body, &requestBody); err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}

		par := formatParam(requestBody)

		logObject := map[string]interface{}{
			"type":            c.Method(),
			"status":          "200",
			"IdUser":          "",
			"IdProject":       idProject,
			"query":           "",
			"requestDuration": 1,
		}

		jsonResponse, err, _ := apiRequest(host+"api/Project/details/"+idProject, "GET", nil)

		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logObject["IdUser"] = jsonResponse["idUser"]

		if key != jsonResponse["key"] {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("INCORRECT KEY PROOJECT")
		}
		endpointsRaw, ok := jsonResponse["endpoints"].([]interface{})
		if !ok {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("CANNOT DESERIALIZE ENDPOINTS")
		}
		//var endpoints []endpoint.Endpoint
		finded := false
		query := ""
		for _, e := range endpointsRaw {
			endpointMap, ok := e.(map[string]interface{})
			if !ok {
				logObject["status"] = "400"
				logRequest(logObject)
				return c.SendString("Invalid endpoint format")
			}

			if url, ok := endpointMap["url"].(string); ok {
				if url == endpointURL {
					finded = true
					query, ok = endpointMap["query"].(string)
					if !ok {
						query = ""
					}
				}
			}
		}

		if !finded {
			logObject["status"] = "404"
			logRequest(logObject)
			return c.SendString("COULD NOT FIND ANY ENDPOINT")
		}
		logObject["query"] = query

		data := map[string]interface{}{
			"database": jsonResponse["iddatabaseNavigation"],
			"query":    query,
			"params":   par,
		}

		dataRequest, err := json.Marshal(data)
		if err != nil {
			logObject["status"] = "500"
			logRequest(logObject)
			return c.SendString("ERROR")
		}

		jsonRes, err, _ := apiRequest(host+"api/Request", "POST", dataRequest)
		//jsonRes, err := http.Post(host+"api/Request", "application/json", bytes.NewBuffer(dataRequest))
		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logRequest(logObject)
		return c.JSON(jsonRes)
	})

	app.Delete("/:key/:projectid/:endpointURL", func(c *fiber.Ctx) error {
		idProject := c.Params("projectid")
		key := c.Params("key")
		endpointURL := c.Params("endpointURL")

		logObject := map[string]interface{}{
			"type":            c.Method(),
			"status":          "200",
			"IdUser":          "",
			"IdProject":       idProject,
			"query":           "",
			"requestDuration": 1,
		}

		m := c.Queries()
		par := formatParam(m)
		jsonResponse, err, _ := apiRequest(host+"api/Project/details/"+idProject, "GET", nil)

		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logObject["IdUser"] = jsonResponse["idUser"]

		if key != jsonResponse["key"] {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("INCORRECT KEY PROOJECT")
		}
		endpointsRaw, ok := jsonResponse["endpoints"].([]interface{})
		if !ok {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("CANNOT DESERIALIZE ENDPOINTS")
		}
		//var endpoints []endpoint.Endpoint
		finded := false
		query := ""
		for _, e := range endpointsRaw {
			endpointMap, ok := e.(map[string]interface{})
			if !ok {
				logObject["status"] = "400"
				logRequest(logObject)
				return c.SendString("Invalid endpoint format")
			}

			if url, ok := endpointMap["url"].(string); ok {
				if url == endpointURL {
					finded = true
					query, ok = endpointMap["query"].(string)
					if !ok {
						query = ""
					}
				}
			}
		}

		if !finded {
			logObject["status"] = "404"
			logRequest(logObject)
			return c.SendString("COULD NOT FIND ANY ENDPOINT")
		}
		logObject["query"] = query

		data := map[string]interface{}{
			"database": jsonResponse["iddatabaseNavigation"],
			"query":    query,
			"params":   par,
		}

		dataRequest, err := json.Marshal(data)
		if err != nil {
			logObject["status"] = "500"
			logRequest(logObject)
			return c.SendString("ERROR")
		}

		jsonRes, err, _ := apiRequest(host+"api/Request", "POST", dataRequest)
		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logRequest(logObject)
		return c.JSON(jsonRes)
	})

	app.Put("/:key/:projectid/:endpointURL", func(c *fiber.Ctx) error {
		idProject := c.Params("projectid")
		key := c.Params("key")
		endpointURL := c.Params("endpointURL")

		logObject := map[string]interface{}{
			"type":            c.Method(),
			"status":          "200",
			"IdUser":          "",
			"IdProject":       idProject,
			"query":           "",
			"requestDuration": 1,
		}

		body := c.Body()
		var requestBody map[string]string

		if err := json.Unmarshal(body, &requestBody); err != nil {
			return c.Status(400).SendString("Invalid JSON")
		}
		par := formatParam(requestBody)
		jsonResponse, err, _ := apiRequest(host+"api/Project/details/"+idProject, "GET", nil)

		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logObject["IdUser"] = jsonResponse["idUser"]

		if key != jsonResponse["key"] {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("INCORRECT KEY PROOJECT")
		}
		endpointsRaw, ok := jsonResponse["endpoints"].([]interface{})
		if !ok {
			logObject["status"] = "400"
			logRequest(logObject)
			return c.SendString("CANNOT DESERIALIZE ENDPOINTS")
		}
		//var endpoints []endpoint.Endpoint
		finded := false
		query := ""
		for _, e := range endpointsRaw {
			endpointMap, ok := e.(map[string]interface{})
			if !ok {
				logObject["status"] = "400"
				logRequest(logObject)
				return c.SendString("Invalid endpoint format")
			}

			if url, ok := endpointMap["url"].(string); ok {
				if url == endpointURL {
					finded = true
					query, ok = endpointMap["query"].(string)
					if !ok {
						query = ""
					}
				}
			}
		}

		if !finded {
			logObject["status"] = "404"
			logRequest(logObject)
			return c.SendString("COULD NOT FIND ANY ENDPOINT")
		}
		logObject["query"] = query

		data := map[string]interface{}{
			"database": jsonResponse["iddatabaseNavigation"],
			"query":    query,
			"params":   par,
		}

		dataRequest, err := json.Marshal(data)
		if err != nil {
			logObject["status"] = "500"
			logRequest(logObject)
			return c.SendString("ERROR")
		}

		jsonRes, err, _ := apiRequest(host+"api/Request", "POST", dataRequest)
		if err != nil {
			logObject["status"] = "500"
			return c.SendString(err.Error())
		}

		logRequest(logObject)
		return c.JSON(jsonRes)
	})

	app.Listen(":3000")

}

func formatParam(m map[string]string) string {
	par := ""
	for key, value := range m {
		if strings.Contains(key, ":string") {
			par += fmt.Sprintf("'%v',", value)
		} else {

			par += fmt.Sprintf("%v,", value)
		}

	}

	return par
}

func apiRequest(url string, method string, body []byte) (map[string]interface{}, error, []map[string]interface{}) {

	var jsonResponse map[string]interface{}
	var mapJsonResponse []map[string]interface{}
	var response *http.Response
	var err error
	user := os.Getenv("API_USER")
	passwd := os.Getenv("API_PASSWORD")
	auth := user + ":" + passwd
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

	var req *http.Request

	if body != nil && len(body) > 0 {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil) // Sin cuerpo
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json") // Asegúrate de establecer el tipo de contenido

	// Realiza la solicitud
	client := &http.Client{}
	response, err = client.Do(req)

	if err != nil {
		return jsonResponse, err, mapJsonResponse
	}

	defer response.Body.Close() // Asegúrate de cerrar el cuerpo de la respuesta

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		return jsonResponse, err, mapJsonResponse
	}

	// Intenta deserializar la respuesta en jsonResponse
	err = json.Unmarshal(resBody, &jsonResponse)
	if err != nil {
		err = json.Unmarshal(resBody, &mapJsonResponse)
		if err != nil {
			return jsonResponse, err, mapJsonResponse
		}
	}

	return jsonResponse, err, mapJsonResponse

}

func logRequest(logObject map[string]interface{}) {
	jsonData, err := json.Marshal(logObject)
	if err != nil {
		return
	}
	resp, err, _ := apiRequest(host+"api/Log", "POST", jsonData)
	//resp, err := http.Post(host+"api/Log", "application/json", nil)

	if err != nil && resp == nil {
		return
	}

}
