package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	//endpoint "github.com/shatigerr/easydev_go/Endpoint"
)

const host = "http://localhost:5102/"

func main() {

	app := fiber.New()

	app.Get("/:key/:projectid/:endpointURL", func(c *fiber.Ctx) error {
		m := c.Queries()
		par := ""
		for key, value := range m {
			if strings.Contains(":string", key) {
				par += fmt.Sprintf("'%v',", value)
			}
			par += fmt.Sprintf("%v,", value)

		}

		idProject := c.Params("projectid")
		key := c.Params("key")
		endpointURL := c.Params("endpointURL")
		response, err := http.Get(host + "api/Project/details/" + idProject)

		if err != nil {
			return c.SendStatus(400)
		}

		body, err := io.ReadAll(response.Body)

		if err != nil {
			return c.SendString("Cannot read the data!")
		}

		var jsonResponse map[string]interface{} // Usamos un mapa genérico para JSON
		err = json.Unmarshal(body, &jsonResponse)

		if err != nil {
			return c.SendString("JSON CANNOT BE DESERIALIZED")
		}

		if key != jsonResponse["key"] {
			return c.SendString("INCORRECT KEY PROOJECT")
		}
		endpointsRaw, ok := jsonResponse["endpoints"].([]interface{})
		if !ok {
			return c.SendString("CANNOT DESERIALIZE ENDPOINTS")
		}
		//var endpoints []endpoint.Endpoint
		finded := false
		query := ""
		for _, e := range endpointsRaw {
			endpointMap, ok := e.(map[string]interface{})
			if !ok {
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
			return c.SendString("COULD NOT FIND ANY ENDPOINT")
		}

		data := map[string]interface{}{
			"database": jsonResponse["iddatabaseNavigation"],
			"query":    query,
			"params":   par,
		}

		dataRequest, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return c.SendString("ERROR")
		}

		res, err := http.Post(host+"api/Request", "application/json", bytes.NewBuffer(dataRequest))

		apiRes, err := io.ReadAll(res.Body)

		if err != nil {
			return c.SendString("Cannot read the data!")
		}

		var jsonRes []map[string]interface{}
		err = json.Unmarshal(apiRes, &jsonRes)

		if err != nil {
			return c.SendString("CANNOT MAKE REQUEST")
		}
		//falta insertar en tabla log
		return c.JSON(jsonRes)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("")
	})

	app.Listen(":3000")

}
