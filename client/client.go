package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// =====================
// Request / Response DTOs
// =====================

type TraceRequest struct {
	From  uint64   `json:"from"`
	To    uint64   `json:"to"`
	Roles []string `json:"roles"`
}

type APIResponse struct {
	Data []SlotData `json:"data"`
}

type SlotData struct {
	Slot         string          `json:"slot"`
	Role         string          `json:"role"`
	Validator    string          `json:"validator"`
	CommitteeID  string          `json:"committeeID"`
	Consensus    []ConsensusStep `json:"consensus"`
	Decideds     json.RawMessage `json:"decideds"` // null or unknown structure
	Pre          []PrePostMsg    `json:"pre"`
	Post         []PrePostMsg    `json:"post"`
	ProposalData string          `json:"proposalData"`
}

type PrePostMsg struct {
	SSVRoot string `json:"ssvRoot"`
	Signer  int    `json:"signer"`
	Time    string `json:"time"`
}

type ConsensusStep struct {
	Proposal     *Proposal     `json:"proposal"`
	Prepares     []Prepare     `json:"prepares"`
	Commits      []Commit      `json:"commits"`
	RoundChanges []RoundChange `json:"roundChanges"`
}

type Proposal struct {
	Round                     int             `json:"round"`
	SSVRoot                   string          `json:"ssvRoot"`
	Leader                    int             `json:"leader"`
	RoundChangeJustifications json.RawMessage `json:"roundChangeJustifications"`
	PrepareJustifications     json.RawMessage `json:"prepareJustifications"`
	Time                      string          `json:"time"`
}

type Prepare struct {
	Round   int    `json:"round"`
	SSVRoot string `json:"ssvRoot"`
	Signer  int    `json:"signer"`
	Time    string `json:"time"`
}

// Shape assumed; your example has commits: null
type Commit struct {
	Round   int    `json:"round"`
	SSVRoot string `json:"ssvRoot"`
	Signer  int    `json:"signer"`
	Time    string `json:"time"`
}

type RoundChange struct {
	Round           int             `json:"round"`
	SSVRoot         string          `json:"ssvRoot"`
	Signer          int             `json:"signer"`
	Time            string          `json:"time"`
	PreparedRound   int             `json:"preparedRound"`
	PrepareMessages json.RawMessage `json:"prepareMessages"`
}

// =====================
// HTTP Client
// =====================

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetValidatorTraces sends:
//
// POST {baseURL}/v1/exporter/traces/validator
// Content-Type: application/json
//
// Body: { "from": ..., "to": ..., "roles": [...] }
func (c *Client) GetValidatorTraces(
	ctx context.Context,
	from, to uint64,
	roles []string,
) (*APIResponse, error) {

	reqPayload := TraceRequest{
		From:  from,
		To:    to,
		Roles: roles,
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	url := c.baseURL + "/v1/exporter/traces/validator"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read a small chunk of body for error context
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}

	var out APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

// =====================
// Example usage
// =====================

func main() {
	ctx := context.Background()

	client := NewClient("http://mainnet-ssv-node-exporter-1.production.vnet.ops.ssvlabsinternal.com")

	resp, err := client.GetValidatorTraces(
		ctx,
		13103596,
		13103696,
		[]string{"PROPOSER"},
	)
	if err != nil {
		log.Fatalf("failed to get validator traces: %v", err)
	}

	for _, d := range resp.Data {
		fmt.Printf("Slot %s | Validator %s | Role %s | Committee %s\n",
			d.Slot, d.Validator, d.Role, d.CommitteeID)

		for stepIdx, step := range d.Consensus {
			fmt.Printf("  Consensus step %d:\n", stepIdx)

			if step.Proposal != nil {
				fmt.Printf("    Proposal: round=%d leader=%d ssvRoot=%s time=%s\n",
					step.Proposal.Round,
					step.Proposal.Leader,
					step.Proposal.SSVRoot,
					step.Proposal.Time,
				)
			}

			for _, p := range step.Prepares {
				fmt.Printf("    Prepare: round=%d signer=%d time=%s\n",
					p.Round, p.Signer, p.Time,
				)
			}

			for _, rc := range step.RoundChanges {
				fmt.Printf("    RoundChange: round=%d signer=%d preparedRound=%d time=%s\n",
					rc.Round, rc.Signer, rc.PreparedRound, rc.Time,
				)
			}
		}
	}
}
