package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/crucial-sa/crux/internal/api"
	"github.com/crucial-sa/crux/internal/auth"
	"github.com/crucial-sa/crux/internal/ui"
)

type Service struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type servicesResponse struct {
	Success bool           `json:"success"`
	Data    []Service      `json:"data,omitempty"`
	Error   *api.ErrorInfo `json:"error,omitempty"`
}

type createServiceRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type apiResponse struct {
	Success bool           `json:"success"`
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

func PromptMissingServiceFields(serviceInfo *Service) error {
	if serviceInfo.Name == "" {
		err := ui.Ask("What should we call your service?", &serviceInfo.Name,
			ui.WithDescription("A unique name to identify this service."),
			ui.WithPlaceholder("my-api"),
			ui.WithValidation(func(s string) error {
				if len(strings.TrimSpace(s)) < 3 {
					return fmt.Errorf("name must be at least 3 characters")
				}
				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("failed to ask service name: %v", err)
		}
	}

	if serviceInfo.Image == "" {
		err := ui.Ask("Which container image should it run?", &serviceInfo.Image,
			ui.WithDescription("An OCI image the platform will boot as a microVM."),
			ui.WithPlaceholder("nginx:latest"),
			ui.WithValidation(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("image is required")
				}
				return nil
			}),
		)
		if err != nil {
			return fmt.Errorf("failed to ask service image: %v", err)
		}
	}

	return nil
}

func CreateService(ctx context.Context, session *auth.Session, serviceInfo *Service) error {
	body, err := json.Marshal(createServiceRequest{
		Name:  serviceInfo.Name,
		Image: serviceInfo.Image,
	})
	if err != nil {
		return fmt.Errorf("failed to encode request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/api/v1/services", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to construct req: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", session.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create service: %v", err)
	}

	defer resp.Body.Close()

	var parsedRes apiResponse

	if err := json.NewDecoder(resp.Body).Decode(&parsedRes); err != nil {
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil || len(bodyBytes) == 0 {
				return fmt.Errorf("server returned non-ok with unparsable or no body: %v, %v", resp.StatusCode, readErr)
			}

			return fmt.Errorf("server returned non-ok response: %v, body: %v", resp.Status, string(bodyBytes))
		}

		return fmt.Errorf("failed to parse json: %v", err)
	}

	if !parsedRes.Success {
		return fmt.Errorf("api error: %v - %v", parsedRes.Error.Code, parsedRes.Error.Message)
	}

	return nil
}

func ConfirmCreate(serviceInfo *Service) bool {
	ui.Summary("Create service", [][2]string{
		{"Name", serviceInfo.Name},
		{"Image", serviceInfo.Image},
	})

	return ui.Confirm("Create this service?")
}
