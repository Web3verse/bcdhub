package elastic

import (
	"fmt"

	"github.com/aopoltorzhicky/bcdhub/internal/models"
	"github.com/tidwall/gjson"
)

func parseStringArray(hit gjson.Result, tag string) []string {
	res := make([]string, 0)
	for _, t := range hit.Get(tag).Array() {
		res = append(res, t.String())
	}
	return res
}

func parseContarctFromHit(hit gjson.Result, c *models.Contract) {
	c.ID = hit.Get("_id").String()
	c.Network = hit.Get("_source.network").String()
	c.Level = hit.Get("_source.level").Int()
	c.Timestamp = hit.Get("_source.timestamp").Time().UTC()
	c.Balance = hit.Get("_source.balance").Int()
	c.Language = hit.Get("_source.language").String()

	c.Tags = parseStringArray(hit, "_source.tags")
	c.Hardcoded = parseStringArray(hit, "_source.hardcoded")
	c.Annotations = parseStringArray(hit, "_source.annotations")
	c.Primitives = parseStringArray(hit, "_source.primitives")
	c.FailStrings = parseStringArray(hit, "_source.fail_strings")
	c.Hash = parseStringArray(hit, "_source.hash")
	c.Entrypoints = parseStringArray(hit, "_source.entrypoints")

	c.Address = hit.Get("_source.address").String()
	c.Manager = hit.Get("_source.manager").String()
	c.Delegate = hit.Get("_source.delegate").String()

	c.ProjectID = hit.Get("_source.project_id").String()
}

func getContractQuery(by map[string]interface{}) map[string]interface{} {
	match := []map[string]interface{}{}
	for k, v := range by {
		match = append(match, map[string]interface{}{
			"match": map[string]interface{}{
				k: v,
			},
		})
	}
	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": match,
			},
		},
	}
}

func (e *Elastic) getContract(q map[string]interface{}) (c models.Contract, err error) {
	res, err := e.query(DocContracts, q)
	if err != nil {
		return
	}
	if res.Get("hits.total.value").Int() != 1 {
		return c, fmt.Errorf("Unknown contract: %v", q)
	}
	hit := res.Get("hits.hits.0")
	parseContarctFromHit(hit, &c)
	return
}

func (e *Elastic) getContracts(q map[string]interface{}) ([]models.Contract, error) {
	res, err := e.query(DocContracts, q)
	if err != nil {
		return nil, err
	}

	contracts := make([]models.Contract, 0)
	arr := res.Get("hits.hits").Array()
	for i := range arr {
		var c models.Contract
		parseContarctFromHit(arr[i], &c)
		contracts = append(contracts, c)
	}
	return contracts, nil
}

// GetContract -
func (e *Elastic) GetContract(by map[string]interface{}) (models.Contract, error) {
	query := getContractQuery(by)
	query["_source"] = map[string]interface{}{
		"excludes": []string{"hash"},
	}
	return e.getContract(query)
}

// GetContractField -
func (e *Elastic) GetContractField(by map[string]interface{}, field string) (interface{}, error) {
	query := getContractQuery(by)
	res, err := e.query(DocContracts, query, field)
	if err != nil {
		return nil, err
	}
	if res.Get("hits.total.value").Int() != 1 {
		return nil, fmt.Errorf("Unknown contract: %v", by)
	}
	return res.Get("hits.hits.0._source").Get(field).Value(), nil
}

// FindProjectContracts -
func (e *Elastic) FindProjectContracts(hash []string, minScore float64) ([]models.Contract, error) {
	if len(hash) != 3 {
		return nil, fmt.Errorf("Length of hash array must be 3")
	}
	query := map[string]interface{}{
		"size": 100,
		"_source": map[string]interface{}{
			"excludes": []string{"hash"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"hash": map[string]interface{}{
								"query":     hash[0],
								"fuzziness": 1,
							},
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"hash": map[string]interface{}{
								"query":     hash[1],
								"fuzziness": 1,
							},
						},
					},
					map[string]interface{}{
						"match": map[string]interface{}{
							"hash": map[string]interface{}{
								"query":     hash[2],
								"fuzziness": 1,
							},
						},
					},
				},
			},
		},
		"sort": map[string]interface{}{
			"timestamp": map[string]interface{}{
				"order": "desc",
			},
		},
	}
	return e.getContracts(query)
}

// SearchByText -
func (e *Elastic) SearchByText(text string) ([]models.Contract, error) {
	query := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"hash"},
		},
		"size": 10,
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": fmt.Sprintf("*%s*", text),
				"fields": []string{
					"address^10", "manager^8", "delegate^6", "tags^4", "hardcoded", "annotations", "fail_strings", "entrypoints",
				},
			},
		},
	}
	return e.getContracts(query)
}

// GetContracts -
func (e *Elastic) GetContracts(q map[string]interface{}) ([]models.Contract, error) {
	query := getContractQuery(q)
	query["size"] = 10000
	return e.getContracts(query)
}