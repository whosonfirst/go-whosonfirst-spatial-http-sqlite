package flags

import (
	"errors"
	"flag"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
)

func ValidateWWWFlags(fs *flag.FlagSet) error {

	enable_www, err := lookup.BoolVar(fs, ENABLE_WWW)

	if err != nil {
		return err
	}

	if !enable_www {
		return nil
	}

	log.Printf("-%s flag is true causing the following flags to also be true: -%s\n", ENABLE_WWW, ENABLE_GEOJSON)

	fs.Set(ENABLE_GEOJSON, "true")
	// fs.Set(ENABLE_PROPERTIES, "true")

	init_lat, err := lookup.Float64Var(fs, INITIAL_LATITUDE)

	if err != nil {
		return err
	}

	if !geo.IsValidLatitude(init_lat) {
		return errors.New("Invalid latitude")
	}

	init_lon, err := lookup.Float64Var(fs, INITIAL_LONGITUDE)

	if err != nil {
		return err
	}

	if !geo.IsValidLongitude(init_lon) {
		return errors.New("Invalid longitude")
	}

	init_zoom, err := lookup.IntVar(fs, INITIAL_ZOOM)

	if err != nil {
		return err
	}

	if init_zoom < 1 {
		return errors.New("Invalid zoom")
	}

	path_flags := []string{
		PATH_PREFIX,
		PATH_API,
		PATH_DATA,
		PATH_PING,
		PATH_PIP,
	}

	for _, fl := range path_flags {

		_, err := lookup.StringVar(fs, fl)

		if err != nil {
			return err
		}
	}

	enable_tangram, err := lookup.BoolVar(fs, ENABLE_TANGRAM)

	if err != nil {
		return err
	}

	if enable_tangram {

		nz_keys := []string{
			NEXTZEN_APIKEY,
			NEXTZEN_STYLE_URL,
			NEXTZEN_TILE_URL,
		}

		for _, k := range nz_keys {

			v, err := lookup.StringVar(fs, k)

			if err != nil {
				return err
			}

			if v == "" {
				log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", k)
			}
		}

	} else {

		v, err := lookup.StringVar(fs, LEAFLET_TILE_URL)

		if err != nil {
			return err
		}

		if v == "" {
			log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", LEAFLET_TILE_URL)
		}
	}

	return nil
}
