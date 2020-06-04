// Package build provides builders for Features protobuf types.
package build

import (
	pb "github.com/bitcdr/protovalid/examples/geo/proto"
)

// Coord builds a protobuf Coordinate.
func Coord(longitude, latitude float64) *pb.Features_Coordinate {
	return &pb.Features_Coordinate{
		Longitude: longitude,
		Latitude:  latitude,
	}
}

// Poi builds a protobuf PoiFeature.
func Poi(name string, longitude, latitude float64) *pb.Features_PoiFeature {
	return &pb.Features_PoiFeature{
		Name:       name,
		Coordinate: Coord(longitude, latitude),
	}
}

// PoiMarker builds a protobuf PoiFeature Marker.
func PoiMarker(color string, size pb.Features_PoiFeature_Marker_Size) *pb.Features_PoiFeature_Marker {
	return &pb.Features_PoiFeature_Marker{
		Color: color,
		Size:  size,
	}
}

// Track builds a protobuf TrackFeature.
func Track(name string, coords []*pb.Features_Coordinate) *pb.Features_TrackFeature {
	return &pb.Features_TrackFeature{
		Name:        name,
		Coordinates: coords,
	}
}

// Stroke builds a protobuf TrackFeature Stroke.
func TrackStroke(color string, width int32) *pb.Features_TrackFeature_Stroke {
	return &pb.Features_TrackFeature_Stroke{
		Color: color,
		Width: width,
	}
}
