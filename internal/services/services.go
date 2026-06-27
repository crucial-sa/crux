package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/crucial-sa/crux/internal/api"
	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/ui"
)

type Service struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type servicesResponse struct {
	Success bool           `json:"success"`
	Data    []Service      `json:"data,omitempty"`
	Error   *api.ErrorInfo `json:"error,omitempty"`
}

var client = &http.Client{}

func GetServices(session *auth.Session) ([]Service, error) {
	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/services", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct req: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch services: %v", err)
	}

	defer resp.Body.Close()

	var parsedRes servicesResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsedRes); err != nil {
		if resp.StatusCode != http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil || len(bodyBytes) == 0 {
				return nil, fmt.Errorf("server returned non-ok with unparsable or no body: %v, %v", resp.StatusCode, err)
			}

			errorText := string(bodyBytes)

			return nil, fmt.Errorf("server return non-ok response: %v, body: %v", resp.Status, errorText)
		}

		return nil, fmt.Errorf("failed to parse json: %v", err)
	}

	if !parsedRes.Success {
		return nil, fmt.Errorf("api error: %v - %v", parsedRes.Error.Code, parsedRes.Error.Message)
	}

	return parsedRes.Data, nil
}

func PrintServicesTable(services []Service) {
	headers := []string{"ID", "SERVICE_NAME", "STATUS"}

	var data [][]string

	for _, service := range services {
		data = append(data, []string{strconv.Itoa(int(service.ID)), service.Name, ui.StatusStyle(service.Status)})
	}

	table := ui.RenderTable(headers, data)
	fmt.Print(table)
}
