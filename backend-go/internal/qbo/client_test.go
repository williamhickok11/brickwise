package qbo

import (
	"encoding/json"
	"testing"
)

func TestParseClassNetIncome(t *testing.T) {
	raw := `{
		"Columns": {
			"Column": [
				{"ColTitle": "", "ColType": "Account"},
				{"ColTitle": "Woodbridge", "ColType": "Money"},
				{"ColTitle": "Oak St", "ColType": "Money"},
				{"ColTitle": "Total", "ColType": "Money"}
			]
		},
		"Rows": {
			"Row": [
				{
					"group": "Income",
					"type": "Section",
					"Rows": {
						"Row": []
					}
				},
				{
					"group": "NetIncome",
					"type": "Section",
					"Summary": {
						"ColData": [
							{"value": "Net Income"},
							{"value": "1240.50"},
							{"value": "-320.00"},
							{"value": "920.50"}
						]
					}
				}
			]
		}
	}`

	var report reportJSON
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	got, err := parseClassNetIncome(report)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if got["Woodbridge"] != 1240.50 {
		t.Errorf("Woodbridge = %v, want 1240.50", got["Woodbridge"])
	}
	if got["Oak St"] != -320.00 {
		t.Errorf("Oak St = %v, want -320.00", got["Oak St"])
	}
	if _, ok := got["Total"]; ok {
		t.Errorf("Total column should be skipped")
	}
}
