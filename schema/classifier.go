package schema

import (
	"fmt"
	"strings"
)

func (s Schema) Classify() error {
	indexes, err := classifyGroupOne(s)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTableClassification, err)
	}

	for i := 0; i < len(indexes); i++ {
		s.Tables[indexes[i]].GroupID = 1
	}

	group := 2
	for {
		indexes, err = classifyGroupsButOne(s, group)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrTableClassification, err)
		}

		if len(indexes) == 0 {
			break
		}

		for i := 0; i < len(indexes); i++ {
			s.Tables[indexes[i]].GroupID = group
		}

		group++
	}

	return nil
}

func (s Schema) GroupedTables() [][]TableSchema {
	group := make([][]TableSchema, 0, len(s.Tables))
	ig := 0

	for {
		tables := make([]TableSchema, 0)
		ig++

		for _, table := range s.Tables {
			if ig == table.GroupID {
				tables = append(tables, table)
			}
		}

		if len(tables) == 0 {
			break
		}

		group = append(group, tables)
	}

	return group
}

func classifyGroupOne(s Schema) ([]int, error) {
	indexes := make([]int, 0, len(s.Tables))

	for i, table := range s.Tables {
		levelOne := false
		referenced := false

		for _, filter := range table.Filters {
			if filterReferenceRegExp.MatchString(filter.Value) {
				referenced = true
			} else {
				levelOne = true
			}
		}

		if len(table.Filters) == 0 || (levelOne && !referenced) {
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 {
		return indexes, fmt.Errorf("couldn't find any level one tables")
	}

	return indexes, nil
}

func classifyGroupsButOne(s Schema, group int) ([]int, error) {
	indexes := make([]int, 0, len(s.Tables))

	for i, table := range s.Tables {
		for _, filter := range table.Filters {
			matches := filterReferenceRegExp.FindAllStringSubmatch(filter.Value, -1)

			if matches != nil {
				refTable := matches[0][1]
				index := findTableByName(s, refTable)

				if index < 0 {
					return nil, fmt.Errorf("%s.%s points to unresolvable reference '%s'", table.Name, filter.Name, matches[0][0])
				}

				if s.Tables[index].GroupID+1 == group {
					indexes = append(indexes, i)
				}
			}
		}
	}

	return indexes, nil
}

func findTableByName(s Schema, tname string) int {
	for i, table := range s.Tables {
		name := strings.ToLower(table.Name)
		if tname == name {
			return i
		}
	}

	return -1
}
