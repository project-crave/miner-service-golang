package repository

import (
	"context"
	"crave/miner/cmd/model"
	database "crave/shared/database"
	craveModel "crave/shared/model"
	"fmt"

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
			SET s:Done
			`
		_, err := tx.Run(ctx, sourceQuery, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			return nil, err
		}

		for _, target := range targets {
			targetQuery := fmt.Sprintf(`
                MERGE (t:Node {name: $targetName})
                WITH t
                MATCH (s:Done {name: $sourceName})
                MERGE (s)-[r:%s]->(t)
                SET r.page = $page,
                    r.context = $context,
					r.appearance = $appearance
            `, page.Name())
			_, err := tx.Run(ctx, targetQuery, map[string]interface{}{
				"targetName": target.Name,
				"context":    target.Context,
				"appearance": target.Appearance,
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
		return err
	}

	return nil
}

func (r *Repository) Remove(name string) error {
	ctx := context.Background()
	session := r.neo4j.NewSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
			MATCH (n:Node {name: $name})
			DETACH DELETE n
		`
		result, err := tx.Run(ctx, query, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, err
		}
		if summary.Counters().NodesDeleted() == 0 {
			return nil, fmt.Errorf("node with name '%s' not found", name)
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}
