package main

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/masoncfrancis/limit/internal/config"
	"github.com/masoncfrancis/limit/internal/database"
	"github.com/valyala/fasthttp"
	"io"
	"log"
	"net/http"
)

func makeRequest(method, url string, body []byte, headers http.Header) ([]byte, int, http.Header, error) {
	// Create a new HTTP request to forward the incoming request
	req, err := http.NewRequest(method, url, io.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, 0, nil, err
	}

	// Copy headers from the original request
	for key, values := range headers {
		for _, value := range values {
			req.Header.Set(key, value)
		}
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Copy response headers
	respHeaders := make(http.Header)
	for key, values := range resp.Header {
		for _, value := range values {
			respHeaders.Set(key, value)
		}
	}

	// Copy response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, nil, err
	}

	return respBody, resp.StatusCode, respHeaders, nil
}

func convertHeader(fasthttpHeader *fasthttp.RequestHeader) http.Header {
	header := make(http.Header)

	fasthttpHeader.VisitAll(func(key, value []byte) {
		header.Set(string(key), string(value))
	})

	return header
}

func main() {

	// TODO add descriptions for all functions

	config.SetupEnvVars() // Load environment variables

	// Create a new Redis client
	dbCtx, rdb := database.CreateDbConn()

	// TODO refactor Go Fiber stuff into it's own file

	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "limit",
	})

	// Set up a route to handle incoming requests
	app.All("/*", func(c *fiber.Ctx) error {

		// Check IP address for previous requests
		ip := c.IP()
		shouldRestrict, err := database.CheckIp(rdb, ip, dbCtx, config.GetTimeWindow(), config.GetRateLimit())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// If the IP address has made more requests than allowed, return HTTP code 429
		if shouldRestrict {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		}

		body, statusCode, headers, err := makeRequest(c.Method(), config.GetForwardUrl()+c.Path(), c.Body(), convertHeader(&c.Request().Header))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		// Set response headers
		for key, values := range headers {
			for _, value := range values {
				c.Set(key, value)
			}
		}

		// Set response status code
		c.Status(statusCode)

		// Send response body
		return c.Send(body)
	})

	// TODO test the server using TLS
	if config.GetUseTls() {
		// Start the server with TLS
		fmt.Printf("Server running on port %s with TLS...", config.GetPort())
		err := app.ListenTLS(":"+config.GetPort(), "./ssl/cert.pem", "./ssl/cert.key")
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
		return
	} else {
		// Start the server without TLS
		fmt.Printf("Server running on port %s without TLS...", config.GetPort())
		err := app.Listen(":" + config.GetPort())
		if err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
		return
	}

}
