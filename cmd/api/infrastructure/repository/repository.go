package repository

import (
	"context"
	"crave/miner/cmd/model"
	database "crave/shared/database"
	craveModel "crave/shared/model"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Repository struct {
	neo4j *database.Neo4jWrapper
}

func NewRepository(neo4j *database.Neo4jWrapper) *Repository {
	return &Repository{neo4j: neo4j}
}

func (r *Repository) Save(name string, page craveModel.Page, targets []model.ParsedTarget) error {
	ctx := context.Background()
	session := r.neo4j.NewSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		sourceQuery := `
			MERGE (s:Node {name: $name})
			SET s:Source
			RETURN s
			`
		_, err := tx.Run(ctx, sourceQuery, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			return nil, err
		}

		for _, target := range targets {
			targetQuery := `
                MERGE (t:Node {name: $targetName})
                SET t:Target
                WITH t
                MATCH (s:Source {name: $sourceName})
                MERGE (s)-[r:REFERENCES]->(t)
                SET r.page = $page,
                    r.context = $context
                RETURN t, r
            `
			_, err := tx.Run(ctx, targetQuery, map[string]interface{}{
				"targetName": target.Name,
				"context":    target.Context,
				"page":       page,
				"sourceName": name,
			})
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		return nil
	}

	return nil
}
