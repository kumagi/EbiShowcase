package savegame

import "encoding/json"

func jsonMarshal(model Model) ([]byte, error) { return json.Marshal(model) }
