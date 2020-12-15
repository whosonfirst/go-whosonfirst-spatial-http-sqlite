# go-whosonfirst-spatial-sqlite

## Important

This is work in progress. It may change, probably has bugs and isn't properly documented yet.

The goal is to have a package that conforms to the [database.SpatialDatabase](https://github.com/whosonfirst/go-whosonfirst-spatial#spatialdatabase) interface using [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) and SQLite's [RTree](https://www.sqlite.org/rtree.html) extension.

Also, this is not as fast as it should be. This is largely with the way WOF records are inflated and passed around in order to support GeoJSON output. There is [an open ticket](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite/issues/2) to address this.

## Databases

This code depends on (4) tables as indexed by the `go-whosonfirst-sqlite-features` package:

* [rtree](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#rtree) - this table is used to perform point-in-polygon spatial queries.
* [spr](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#spr) - this table is used to generate [standard place response](#) (SPR) results.
* [geometry](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#geometry) - this table is used to append geometries to GeoJSON-formatted results.
* [properties](https://github.com/whosonfirst/go-whosonfirst-sqlite-features#properties) - this table is used to append extra properties (to the SPR response) for GeoJSON-formatted results.

The `go-whosonfirst-sqlite-features` package also indexes a `geojson` table  but it turns out that retrieving, and parsing, properties and geometries from their own tables is faster.

Here's an example of the creating a compatible SQLite database for all the [administative data in Canada](https://github.com/whosonfirst-data/whosonfirst-data-admin-ca) using the `wof-sqlite-index-features` tool which is part of the [go-whosonfirst-sqlite-features-index](https://github.com/whosonfirst/go-whosonfirst-sqlite-features-index) package:

```
$> ./bin/wof-sqlite-index-features \
	-index-alt-files \
	-rtree \
	-spr \
	-geometry \
	-properties \
	-timings \
	-dsn /usr/local/ca-alt.db \
	-mode repo:// \
	/usr/local/data/whosonfirst-data-admin-ca/

13:09:44.642004 [wof-sqlite-index-features] STATUS time to index rtree (11860) : 30.469010289s
13:09:44.642136 [wof-sqlite-index-features] STATUS time to index geometry (11860) : 5.155172377s
13:09:44.642141 [wof-sqlite-index-features] STATUS time to index properties (11860) : 4.631908497s
13:09:44.642143 [wof-sqlite-index-features] STATUS time to index spr (11860) : 19.160260741s
13:09:44.642146 [wof-sqlite-index-features] STATUS time to index all (11860) : 1m0.000182571s
13:10:44.642848 [wof-sqlite-index-features] STATUS time to index spr (32724) : 39.852608874s
13:10:44.642861 [wof-sqlite-index-features] STATUS time to index rtree (32724) : 57.361318918s
13:10:44.642864 [wof-sqlite-index-features] STATUS time to index geometry (32724) : 10.242155898s
13:10:44.642868 [wof-sqlite-index-features] STATUS time to index properties (32724) : 10.815961878s
13:10:44.642871 [wof-sqlite-index-features] STATUS time to index all (32724) : 2m0.000429956s
```

And then...

```
$> ./bin/query \
	-database-uri 'sqlite://?dsn=/usr/local/data/ca-alt.db' \
	-latitude 45.572744 \
	-longitude -73.586295
| jq \
| grep wof:id

2020/12/15 15:32:05 Unable to parse placetype (alt) for ID 85874359, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Unable to parse placetype (alt) for ID 85633041, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Unable to parse placetype (alt) for ID 136251273, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Unable to parse placetype (alt) for ID 85633041, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Unable to parse placetype (alt) for ID 136251273, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Unable to parse placetype (alt) for ID 85633041, because 'Invalid placetype' - skipping placetype filters
2020/12/15 15:32:06 Time to point in polygon, 596.579126ms
      "wof:id": "1108955735",
      "wof:id": "85874359",
      "wof:id": "85874359",
      "wof:id": "890458661",
      "wof:id": "85633041",
      "wof:id": "136251273",
      "wof:id": "85633041",
      "wof:id": "85633041",
      "wof:id": "136251273",
      "wof:id": "85633041",
      "wof:id": "136251273",
```

_TBW: Indexing tables on start-up._

## Example

```
package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

func main() {

	database_uri := "sqlite://?dsn=whosonfirst.db"
	properties_uri := "sqlite://?dsn=whosonfirst.db"
	latitude := 37.616951
	longitude := -122.383747

	props := []string{
		"wof:concordances",
		"wof:hierarchy",
		"sfomuseum:*",
	}

	ctx := context.Background()
	
	db, _ := database.NewSpatialDatabase(ctx, *database_uri)
	pr, _ := properties.NewPropertiesReader(ctx, *properties_uri)
	
	c, _ := geo.NewCoordinate(*longitude, *latitude)
	f, _ := filter.NewSPRFilter()
	r, _ := db.PointInPolygon(ctx, c, f)

	r, _ = pr.PropertiesResponseResultsWithStandardPlacesResults(ctx, r, props)

	enc, _ := json.Marshal(r)
	fmt.Println(string(enc))
}
```

_Error handling removed for the sake of brevity._

## Interfaces

This package implements the following [go-whosonfirst-spatial](#) interfaces.

### spatial.SpatialDatabase

```
import (
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"       
)

db, err := database.NewSpatialDatabase(ctx, "sqlite://?dsn={DSN}")
```

### spatial.PropertiesReader

```
import (
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"       
)

pr, err := properties.NewPropertiesReader(ctx, "sqlite://?dsn={DSN}")
```

## Tools

### query

```
$> ./bin/query -h
Usage of ./bin/query:
  -database-uri string
    	...
  -latitude float
    	...
  -longitude float
    	...
  -properties value
    	...
  -properties-uri string
    	...
```

For example:

```
$> ./bin/query \
	-database-uri 'sqlite://?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-properties-uri 'sqlite://?dsn=/usr/local/data/sfomuseum-data-architecture.db' \
	-latitude 37.616951 \
	-longitude -122.383747 \
	-properties 'wof:hierarchy' \
	-properties 'sfomuseum:*' \
| jq

{
  "properties": [
    {
      "mz:is_ceased": 1,
      "mz:is_current": 0,
      "mz:is_deprecated": 0,
      "mz:is_superseded": 1,
      "mz:is_superseding": 1,
      "mz:latitude": 37.617475,
      "mz:longitude": -122.383371,
      "mz:max_latitude": 37.61950174060331,
      "mz:max_longitude": -122.38139655218178,
      "mz:min_latitude": 37.61615511156664,
      "mz:min_longitude": -122.3853565208227,
      "mz:uri": "https://data.whosonfirst.org/115/939/616/5/1159396165.geojson",
      "sfomuseum:is_sfo": 1,
      "sfomuseum:placetype": "terminal",
      "sfomuseum:terminal_id": "CENTRAL",
      "wof:country": "US",
      "wof:hierarchy": [
        {
          "building_id": 1159396339,
          "campus_id": 102527513,
          "continent_id": 102191575,
          "country_id": 85633793,
          "county_id": 102087579,
          "locality_id": 85922583,
          "neighbourhood_id": -1,
          "region_id": 85688637,
          "wing_id": 1159396165
        }
      ],
      "wof:id": 1159396165,
      "wof:lastmodified": 1547232162,
      "wof:name": "Central Terminal",
      "wof:parent_id": 1159396339,
      "wof:path": "115/939/616/5/1159396165.geojson",
      "wof:placetype": "wing",
      "wof:repo": "sfomuseum-data-architecture",
      "wof:superseded_by": [
        1159396149
      ],
      "wof:supersedes": [
        1159396171
      ]
    },

    ... and so on
   }
]   
```

Note: This assumes a database that was previously indexed using the [whosonfirst/go-whosonfirst-sqlite-features](https://github.com/whosonfirst/go-whosonfirst-sqlite-features) `wof-sqlite-index-features` tool. For example:

```
$> ./bin/wof-sqlite-index-features -rtree -geojson -dsn /tmp/test.db -mode repo:// /usr/local/data/sfomuseum-data-architecture/
```

## See also

* https://www.sqlite.org/rtree.html
* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-sqlite
* https://github.com/whosonfirst/go-whosonfirst-sqlite-features