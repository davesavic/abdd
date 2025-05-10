package app

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func (a *Abdd) GenerateFakeData(t *Test) error {
	fmt.Println("Generating fake data...")

	if t.Fake == nil {
		return nil
	}

	for key, tag := range t.Fake {
		value, err := generateFromTag(tag)
		if err != nil {
			return fmt.Errorf("failed to generate fake data for %s: %w", key, err)
		}
		a.Store[key] = value
	}

	return nil
}

func generateFromTag(tag string) (string, error) {
	faker := gofakeit.New(0)
	result, err := faker.Generate(tag)
	if err != nil {
		return "", fmt.Errorf("failed to generate data from tag: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
}

// func (a *Abdd)
