package services

import (
	"fmt"
	"strconv"

	"github.com/yapi-teklif/internal/forms"
	"github.com/yapi-teklif/internal/repositories"
)

type IConstructionFeatureService interface {
	CheckFeaturesTypes(form []forms.ConstructionFeatureForm, categoryID uint) error
}

type ConstructionFeatureService struct {
	Repository repositories.IConstructionFeatureRepository
}

func GetConstructionFeatureService(r repositories.IConstructionFeatureRepository) *ConstructionFeatureService {
	return &ConstructionFeatureService{Repository: r}
}

func (s *ConstructionFeatureService) CheckFeaturesTypes(form []forms.ConstructionFeatureForm, categoryID uint) error {
	// Retrieve all features belonging to the product_category
	features, err := s.Repository.GetFeaturesByCategoryID(categoryID)
	if err != nil {
		return err
	}

	// Create a map of feature IDs to track which features have been provided in the form
	providedFeatures := make(map[uint]bool)
	for _, v := range form {
		providedFeatures[v.FeatureID] = true
	}

	// Check if all required features are present in the form
	for _, feature := range features {
		if feature.IsRequired && !providedFeatures[feature.ID] {
			return fmt.Errorf("required feature with ID %d is not present in the form", feature.ID)
		}
	}

	// Process each form element sequentially
	for _, v := range form {
		feature, err := s.Repository.GetFeatureTypeById(v.FeatureID)
		if err != nil {
			return err
		}
		featureType := feature.Type

		switch featureType {
		case "integer":
			_, err := strconv.Atoi(v.Value)
			if err != nil {
				return fmt.Errorf("invalid value for feature with ID %d: %w", v.FeatureID, err)
			}
		case "float64":
			_, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				return fmt.Errorf("invalid value for feature with ID %d: %w", v.FeatureID, err)
			}
		case "boolean":
			_, err := strconv.ParseBool(v.Value)
			if err != nil {
				return fmt.Errorf("invalid value for feature with ID %d: %w", v.FeatureID, err)
			}
		case "string":
			// No need to check for errors here
		}
	}

	return nil
}
