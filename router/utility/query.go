package utility

import (
	"fmt"
	"server/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/mitchellh/mapstructure"
)

type validOps struct {
	EQ    interface{} `json:"eq,omitempty"`
	GT    interface{} `json:"gt,omitempty"`
	GTE   interface{} `json:"gte,omitempty"`
	LT    interface{} `json:"lt,omitempty"`
	LTE   interface{} `json:"lte,omitempty"`
	ILIKE interface{} `json:"ilike,omitempty"`
	LIKE  interface{} `json:"like,omitempty"`
}

func determineValidOps(template sq.SelectBuilder, key string, value interface{}) sq.SelectBuilder {
	var ops validOps
	err := mapstructure.Decode(value, &ops)
	if err != nil {
		return template
	}
	if ops.EQ != nil {
		template = template.Where(sq.Eq{fmt.Sprintf(`"%s"`, key): ops.EQ})
		return template
	}
	if ops.GT != nil {
		template = template.Where(sq.Gt{fmt.Sprintf(`"%s"`, key): ops.GT})
		return template
	}
	if ops.GTE != nil {
		template = template.Where(sq.GtOrEq{fmt.Sprintf(`"%s"`, key): ops.GTE})
		return template
	}
	if ops.LT != nil {
		template = template.Where(sq.Lt{fmt.Sprintf(`"%s"`, key): ops.LT})
		return template
	}
	if ops.LTE != nil {
		template = template.Where(sq.LtOrEq{fmt.Sprintf(`"%s"`, key): ops.LTE})
		return template
	}
	if ops.ILIKE != nil {
		template = template.Where(sq.ILike{fmt.Sprintf(`"%s"`, key): fmt.Sprint("%", ops.ILIKE, "%")})
		return template
	}
	if ops.LIKE != nil {
		template = template.Where(sq.Like{fmt.Sprintf(`"%s"`, key): fmt.Sprint("%", ops.LIKE, "%")})
		return template
	}
	return template
}

func checkInFilters(value string, validFilters []string) bool {
	for i := 0; i < len(validFilters); i++ {
		if value == validFilters[i] {
			return true
		}
	}
	return false
}

func BuildWhere(template sq.SelectBuilder, validFilters []string, query models.Filter) sq.SelectBuilder {
	if query == nil {
		return template
	}

	for key, value := range query {
		if !checkInFilters(key, validFilters) {
			continue
		}
		switch value.(type) {
		case string:
			template = template.Where(fmt.Sprintf(`"%s"`, key)+" ILIKE ? ", fmt.Sprint("%", value, "%"))
		case int:
			template = template.Where(fmt.Sprintf(`"%s"`, key)+"= ? ", value)
		}
		template = determineValidOps(template, key, value)
	}

	return template
}
