package trello

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	apiKey     string
	token      string
	baseURL    string
	httpClient *http.Client
}

type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Closed bool   `json:"closed"`
	URL    string `json:"url"`
}

type List struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Closed bool   `json:"closed"`
}

type Card struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Desc             string   `json:"desc"`
	Due              string   `json:"due"`
	URL              string   `json:"url"`
	IDList           string   `json:"idList"`
	DateLastActivity string   `json:"dateLastActivity"`
	IDLabels         []string `json:"idLabels"`
}

type Label struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	IDBoard string `json:"idBoard"`
}

func NewClient(apiKey, token string) *Client {
	return &Client{
		apiKey:     apiKey,
		token:      token,
		baseURL:    "https://api.trello.com/1",
		httpClient: &http.Client{},
	}
}

func (c *Client) GetBoards() ([]Board, error) {
	params := c.authParams()

	endpoint := c.baseURL +
		"/members/me/boards?" +
		params.Encode()

	var boards []Board

	if err := c.doJSON(
		http.MethodGet,
		endpoint,
		nil,
		&boards,
		http.StatusOK,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch boards: %w", err)
	}

	return boards, nil
}

func (c *Client) GetLists(boardID string) ([]List, error) {
	params := c.authParams()

	endpoint := c.baseURL +
		"/boards/" +
		url.PathEscape(boardID) +
		"/lists?" +
		params.Encode()

	var lists []List

	if err := c.doJSON(
		http.MethodGet,
		endpoint,
		nil,
		&lists,
		http.StatusOK,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch lists: %w", err)
	}

	return lists, nil
}

func (c *Client) GetCards(listID string) ([]Card, error) {
	params := c.authParams()

	params.Set(
		"fields",
		"id,name,desc,due,url,idList,dateLastActivity,idLabels",
	)

	endpoint := c.baseURL +
		"/lists/" +
		url.PathEscape(listID) +
		"/cards?" +
		params.Encode()

	var cards []Card

	if err := c.doJSON(
		http.MethodGet,
		endpoint,
		nil,
		&cards,
		http.StatusOK,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch cards: %w", err)
	}

	return cards, nil
}

func (c *Client) GetLabels(boardID string) ([]Label, error) {
	params := c.authParams()

	endpoint := c.baseURL +
		"/boards/" +
		url.PathEscape(boardID) +
		"/labels?" +
		params.Encode()

	var labels []Label

	if err := c.doJSON(
		http.MethodGet,
		endpoint,
		nil,
		&labels,
		http.StatusOK,
	); err != nil {
		return nil, fmt.Errorf("failed to fetch labels: %w", err)
	}

	return labels, nil
}

func (c *Client) CreateLabel(
	boardID string,
	name string,
	color string,
) (*Label, error) {
	params := c.authParams()

	params.Set("name", name)
	params.Set("color", color)

	endpoint := c.baseURL +
		"/boards/" +
		url.PathEscape(boardID) +
		"/labels?" +
		params.Encode()

	var label Label

	if err := c.doJSON(
		http.MethodPost,
		endpoint,
		nil,
		&label,
		http.StatusOK,
	); err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	return &label, nil
}

func (c *Client) AddLabelToCard(
	cardID string,
	labelID string,
) error {
	params := c.authParams()
	params.Set("value", labelID)

	endpoint := c.baseURL +
		"/cards/" +
		url.PathEscape(cardID) +
		"/idLabels?" +
		params.Encode()

	if err := c.doJSON(
		http.MethodPost,
		endpoint,
		nil,
		nil,
		http.StatusOK,
	); err != nil {
		return fmt.Errorf(
			"failed to add label to card: %w",
			err,
		)
	}

	return nil
}

func (c *Client) authParams() url.Values {
	params := url.Values{}

	params.Set("key", c.apiKey)
	params.Set("token", c.token)

	return params
}

func (c *Client) doJSON(
	method string,
	endpoint string,
	body io.Reader,
	result any,
	expectedStatus int,
) error {
	req, err := http.NewRequest(
		method,
		endpoint,
		body,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set(
			"Content-Type",
			"application/json",
		)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		responseBody, _ := io.ReadAll(
			io.LimitReader(resp.Body, 4096),
		)

		message := strings.TrimSpace(
			string(responseBody),
		)

		if message != "" {
			return fmt.Errorf(
				"trello returned %s: %s",
				resp.Status,
				message,
			)
		}

		return fmt.Errorf(
			"trello returned status: %s",
			resp.Status,
		)
	}

	if result == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf(
			"failed to decode response: %w",
			err,
		)
	}

	return nil
}
