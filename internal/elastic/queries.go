package elastic

var allFields = []string{
	"address^10", "manager^8", "delegate^6", "tags^4", "hardcoded", "annotations", "fail_strings", "entrypoints", "language",
}

var mapFields = map[string]string{
	"address":   "address^10",
	"manager":   "manager^8",
	"delegate":  "delegate^6",
	"tags":      "tags^4",
	"hardcoded": "hardcoded",
	"annots":    "annotations",
	"fail":      "fail_strings",
	"entry":     "entrypoints",
	"language":  "language",
}

var mapHighlights = qItem{
	"address":      qItem{},
	"manager":      qItem{},
	"delegate":     qItem{},
	"tags":         qItem{},
	"hardcoded":    qItem{},
	"annotations":  qItem{},
	"fail_strings": qItem{},
	"entrypoints":  qItem{},
	"language":     qItem{},
}

var supportedNetworks = map[string]struct{}{
	"mainnet":     struct{}{},
	"zeronet":     struct{}{},
	"babylonnet":  struct{}{},
	"carthagenet": struct{}{},
}
