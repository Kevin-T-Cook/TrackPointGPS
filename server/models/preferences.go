package models

import "encoding/json"

type Preferences struct {
    ID            uint            `json:"id"`
    UserID        string          `json:"user_id"`
    SortOrder     string          `json:"sort_order"`
    HiddenDevices json.RawMessage `json:"hidden_devices"`
}