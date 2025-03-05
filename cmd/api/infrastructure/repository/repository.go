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

func (r *Repository) SaveOrigin(name string, tag int64) error {
	ctx := context.Background()
	session := r.neo4j.NewSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {

		tagQuery := fmt.Sprintf(`
			MERGE (s:Node {name: $name})
			RETURN s.tag
			`)
		result, err := tx.Run(ctx, tagQuery, map[string]interface{}{
			"name": name,
		})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			tagValue := result.Record().Values[0]
			if intValue, ok := tagValue.(int64); ok {
				tag |= intValue
			}
		}

		query := `
		MERGE (n:Node {name: $name})
		SET n:Done,
		n.tag = $newTag
		`

		params := map[string]interface{}{
			"name":   name,
			"newTag": tag,
		}

		result, err = tx.Run(ctx, query, params)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}

		if _, err := result.Consume(ctx); err != nil {
			return nil, fmt.Errorf("failed to consume result: %w", err)
		}
		if err != nil {
			return nil, err
		}
		return nil, err
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SaveDestination(org string, dest *model.ParsedTarget, page craveModel.Page, tag int64) error {
	ctx := context.Background()
	session := r.neo4j.NewSession(ctx)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {

		tagQuery := `
			MERGE (s:Node {name: $name})
			RETURN s.tag
			`
		result, err := tx.Run(ctx, tagQuery, map[string]interface{}{
			"name": dest.Name,
		})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			tagValue := result.Record().Values[0]
			if intValue, ok := tagValue.(int64); ok {
				tag &= intValue
			}
		}
		destQuery := fmt.Sprintf(`
                MERGE (d:Node {name: $destName})
                WITH d
                MATCH (o:Done {name: $orgName})
                MERGE (o)-[r:%s]->(d)
				SET 
                    r.context = $context,
					r.appearance = $appearance
				SET
					d.tag = $newTag
            `, page.Name())

		result, err = tx.Run(ctx, destQuery, map[string]interface{}{
			"destName":   dest.Name,
			"orgName":    org,
			"context":    dest.Context,
			"appearance": dest.Appearance,
			"newTag":     tag,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		if _, err = result.Consume(ctx); err != nil {
			return nil, fmt.Errorf("failed to consume result: %w", err)
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
