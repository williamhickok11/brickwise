package qbo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"brickwise/backend/internal/models"
)

const (
	oauthAuthorizeURL = "https://appcenter.intuit.com/connect/oauth2"
	oauthTokenURL     = "https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer"
	sandboxAPIBase    = "https://sandbox-quickbooks.api.intuit.com"
	productionAPIBase = "https://quickbooks.api.intuit.com"
	defaultScope      = "com.intuit.quickbooks.accounting"
)

// Config holds Intuit app credentials from the environment.
type Config struct {
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	Environment      string // sandbox | production
	FrontendRedirect string
}

// LoadConfig reads QBO_* environment variables.
func LoadConfig() Config {
	redirect := os.Getenv("QBO_REDIRECT_URI")
	if redirect == "" {
		redirect = "http://localhost:8080/api/v1/integrations/qbo/callback"
	}
	frontend := os.Getenv("QBO_FRONTEND_REDIRECT")
	if frontend == "" {
		frontend = "http://localhost:5173/settings/quickbooks"
	}
	env := os.Getenv("QBO_ENVIRONMENT")
	if env == "" {
		env = "sandbox"
	}
	return Config{
		ClientID:         os.Getenv("QBO_CLIENT_ID"),
		ClientSecret:     os.Getenv("QBO_CLIENT_SECRET"),
		RedirectURI:      redirect,
		Environment:      env,
		FrontendRedirect: frontend,
	}
}

// Configured reports whether OAuth credentials are present.
func (c Config) Configured() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

func (c Config) apiBase() string {
	if strings.EqualFold(c.Environment, "production") {
		return productionAPIBase
	}
	return sandboxAPIBase
}

// Client talks to Intuit OAuth and the QBO Accounting API.
type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
	}
}

// AuthorizeURL builds the Intuit OAuth consent URL.
func (c *Client) AuthorizeURL(state string) string {
	q := url.Values{}
	q.Set("client_id", c.cfg.ClientID)
	q.Set("redirect_uri", c.cfg.RedirectURI)
	q.Set("response_type", "code")
	q.Set("scope", defaultScope)
	q.Set("state", state)
	return oauthAuthorizeURL + "?" + q.Encode()
}

// TokenResponse is the OAuth token endpoint payload.
type TokenResponse struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	ExpiresIn              int    `json:"expires_in"`
	XRefreshTokenExpiresIn int    `json:"x_refresh_token_expires_in"`
	TokenType              string `json:"token_type"`
}

// ExchangeCode swaps an authorization code for tokens.
func (c *Client) ExchangeCode(ctx context.Context, code string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", c.cfg.RedirectURI)
	return c.postToken(ctx, form)
}

// RefreshAccessToken refreshes an expired access token.
func (c *Client) RefreshAccessToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)
	return c.postToken(ctx, form)
}

func (c *Client) postToken(ctx context.Context, form url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(c.cfg.ClientID, c.cfg.ClientSecret))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, truncate(string(body), 300))
	}

	var tokens TokenResponse
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	return &tokens, nil
}

// ListClasses returns active Classes from the company file.
func (c *Client) ListClasses(ctx context.Context, realmID, accessToken string) ([]models.QBOClass, error) {
	query := "SELECT Id, Name FROM Class MAXRESULTS 1000"
	endpoint := fmt.Sprintf("%s/v3/company/%s/query?query=%s&minorversion=73",
		c.cfg.apiBase(), url.PathEscape(realmID), url.QueryEscape(query))

	var payload struct {
		QueryResponse struct {
			Class []struct {
				ID   string `json:"Id"`
				Name string `json:"Name"`
			} `json:"Class"`
		} `json:"QueryResponse"`
	}
	if err := c.getJSON(ctx, endpoint, accessToken, &payload); err != nil {
		return nil, err
	}

	out := make([]models.QBOClass, 0, len(payload.QueryResponse.Class))
	for _, cl := range payload.QueryResponse.Class {
		out = append(out, models.QBOClass{ID: cl.ID, Name: cl.Name})
	}
	return out, nil
}

// ClassNetIncome maps Class name → cash-basis net income for a period.
func (c *Client) ClassNetIncome(ctx context.Context, realmID, accessToken string, start, end time.Time) (map[string]float64, error) {
	q := url.Values{}
	q.Set("start_date", start.Format("2006-01-02"))
	q.Set("end_date", end.Format("2006-01-02"))
	q.Set("accounting_method", "Cash")
	q.Set("summarize_column_by", "Class")
	q.Set("minorversion", "73")

	endpoint := fmt.Sprintf("%s/v3/company/%s/reports/ProfitAndLoss?%s",
		c.cfg.apiBase(), url.PathEscape(realmID), q.Encode())

	var report reportJSON
	if err := c.getJSON(ctx, endpoint, accessToken, &report); err != nil {
		return nil, err
	}
	return parseClassNetIncome(report)
}

func (c *Client) getJSON(ctx context.Context, endpoint, accessToken string, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qbo request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf("qbo api returned %d: %s", resp.StatusCode, truncate(string(body), 400))
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("decode qbo response: %w", err)
	}
	return nil
}

type reportJSON struct {
	Columns struct {
		Column []struct {
			ColTitle string `json:"ColTitle"`
			ColType  string `json:"ColType"`
			MetaData []struct {
				Name  string `json:"Name"`
				Value string `json:"Value"`
			} `json:"MetaData"`
		} `json:"Column"`
	} `json:"Columns"`
	Rows reportRows `json:"Rows"`
}

type reportRows struct {
	Row []reportRow `json:"Row"`
}

type reportRow struct {
	Header  *reportSection `json:"Header,omitempty"`
	Rows    *reportRows    `json:"Rows,omitempty"`
	Summary *reportSection `json:"Summary,omitempty"`
	ColData []colData      `json:"ColData,omitempty"`
	Type    string         `json:"type"`
	Group   string         `json:"group"`
}

type reportSection struct {
	ColData []colData `json:"ColData"`
}

type colData struct {
	Value string `json:"value"`
	ID    string `json:"id,omitempty"`
}

func parseClassNetIncome(report reportJSON) (map[string]float64, error) {
	classNames := make([]string, 0)
	for i, col := range report.Columns.Column {
		if i == 0 {
			continue
		}
		title := strings.TrimSpace(col.ColTitle)
		if title == "" || strings.EqualFold(title, "Total") {
			classNames = append(classNames, "")
			continue
		}
		classNames = append(classNames, title)
	}

	netCols := findNetIncomeColumns(report.Rows.Row)
	if netCols == nil {
		return map[string]float64{}, nil
	}

	result := make(map[string]float64)
	for i, name := range classNames {
		if name == "" {
			continue
		}
		idx := i + 1
		if idx >= len(netCols) {
			break
		}
		val, err := parseAmount(netCols[idx].Value)
		if err != nil {
			continue
		}
		result[name] = val
	}
	return result, nil
}

func findNetIncomeColumns(rows []reportRow) []colData {
	for i := range rows {
		row := &rows[i]
		if strings.EqualFold(row.Group, "NetIncome") {
			if row.Summary != nil && len(row.Summary.ColData) > 0 {
				return row.Summary.ColData
			}
			if len(row.ColData) > 0 {
				return row.ColData
			}
		}
		if row.Summary != nil && len(row.Summary.ColData) > 0 {
			if strings.EqualFold(strings.TrimSpace(row.Summary.ColData[0].Value), "Net Income") {
				return row.Summary.ColData
			}
		}
		if row.Rows != nil {
			if found := findNetIncomeColumns(row.Rows.Row); found != nil {
				return found
			}
		}
	}
	return nil
}

func parseAmount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}
	s = strings.ReplaceAll(s, ",", "")
	return strconv.ParseFloat(s, 64)
}

func basicAuth(id, secret string) string {
	return base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// ApplyTokens fills connection fields from a token response.
func ApplyTokens(conn *models.QBOConnection, tokens *TokenResponse, now time.Time) {
	conn.AccessToken = tokens.AccessToken
	conn.RefreshToken = tokens.RefreshToken
	conn.AccessTokenExpiresAt = now.Add(time.Duration(tokens.ExpiresIn) * time.Second)
	if tokens.XRefreshTokenExpiresIn > 0 {
		conn.RefreshTokenExpiresAt = now.Add(time.Duration(tokens.XRefreshTokenExpiresIn) * time.Second)
	} else {
		conn.RefreshTokenExpiresAt = now.Add(100 * 24 * time.Hour)
	}
	conn.UpdatedAt = now
}
