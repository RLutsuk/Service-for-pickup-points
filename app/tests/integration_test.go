package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func Test_IntegrationCreateReception(t *testing.T) {
	baseURL := url.URL{Scheme: "http", Host: host}
	e := httpexpect.Default(t, baseURL.String())

	modToken := getToken(t, e, "moderator")
	empployeeToken := getToken(t, e, "employee")

	ppID := createPickupPoint(t, e, modToken)
	createReception(t, e, empployeeToken, ppID)
	addProducts(t, e, empployeeToken, ppID, 50)
	closeReception(t, e, empployeeToken, ppID)
}

func getToken(t *testing.T, e *httpexpect.Expect, role string) string {
	resp := e.POST("/dummyLogin").WithJSON(map[string]interface{}{
		"role": role,
	}).Expect().Status(http.StatusOK).JSON().Object()

	token := resp.Value("token").String().Raw()
	if token == "" {
		t.Error("expected token, got empty")
	}
	return token
}

func createPickupPoint(t *testing.T, e *httpexpect.Expect, token string) string {
	resp := e.POST("/pvz").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"city": "Москва",
		}).
		Expect().Status(http.StatusCreated).JSON().Object()

	id := resp.Value("id").String().Raw()
	city := resp.Value("city").String().Raw()

	if city != "Москва" || id == "" {
		t.Error("invalid json response, createPickupPoint()")
	}

	return id
}

func createReception(t *testing.T, e *httpexpect.Expect, token, ppID string) {
	resp := e.POST("/receptions").
		WithHeader("Authorization", "Bearer "+token).
		WithJSON(map[string]interface{}{
			"pvzId": ppID,
		}).
		Expect().Status(http.StatusCreated).JSON().Object()

	id := resp.Value("id").String().Raw()
	status := resp.Value("status").String().Raw()
	newppID := resp.Value("pvzId").String().Raw()

	if status != "in_progress" || newppID != ppID || id == "" {
		t.Error("invalid json response, createReception()")
	}
}

func addProducts(t *testing.T, e *httpexpect.Expect, token, ppID string, count int) {
	for i := 0; i < count; i++ {
		e.POST("/products").
			WithHeader("Authorization", "Bearer "+token).
			WithJSON(map[string]interface{}{
				"type":  "электроника",
				"pvzId": ppID,
			}).
			Expect().Status(http.StatusCreated)
	}
}
func closeReception(t *testing.T, e *httpexpect.Expect, token, ppID string) {
	resp := e.POST(fmt.Sprintf("/pvz/%s/close_last_reception", ppID)).
		WithHeader("Authorization", "Bearer "+token).
		Expect().Status(http.StatusOK).JSON().Object()

	status := resp.Value("status").String().Raw()
	newppID := resp.Value("pvzId").String().Raw()

	if status != "close" || newppID != ppID {
		t.Error("invalid json response, createReception()")
	}
}
