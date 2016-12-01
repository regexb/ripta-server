package geocode

import (
	"fmt"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

// Service for geocoding an address or cross streets
type Service interface {
	Geocode(ctx context.Context, request Request) (*Response, error)
}

// NewGeocodeService returns a new instance of a geocode service
func NewGeocodeService(client maps.Client) Service {
	return &geoService{
		client: client,
	}
}

type geoService struct {
	client maps.Client
}

func (g *geoService) Geocode(ctx context.Context, request Request) (*Response, error) {
	// setup components to scope geocode requests
	components := make(map[maps.Component]string)
	components[maps.ComponentAdministrativeArea] = "Rhode Island"
	components[maps.ComponentCountry] = "us"

	geoRequest := &maps.GeocodingRequest{
		Address:    request.Query,
		Components: components,
	}
	addresses, err := g.client.Geocode(ctx, geoRequest)
	if err != nil {
		return nil, fmt.Errorf("Failed to lookup lat and long for %s: Error %v", request.Query, err)
	}

	geocodedAddresses := []*GeocodeAddress{}
	for _, a := range addresses {
		geocodedAddresses = append(geocodedAddresses, &GeocodeAddress{
			Address: a.FormattedAddress,
			Lat:     a.Geometry.Location.Lat,
			Long:    a.Geometry.Location.Lng,
		})
	}

	return &Response{geocodedAddresses}, nil
}
