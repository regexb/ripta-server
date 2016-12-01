# RIPTA Api Server
-----------------

## Endpoints

- [X] *route_stops*: `/api/v1/routes/<route_id>?stop_id=<stop_id>`
- [X] *stop_locations*: `/api/v1/stop?location=<lat,long>`
- [X] *geocode*: `/api/v1/geocode?q=<query_string>`

#### Route Stops

Returns the current list of upcoming bus stops for a given route. Can filter by stop_id

#### Stop Locations

Returns a list of the closest bus stops to the provided lat and long coordinates.

#### Geocode

Uses an address or cross street to look up the lat and long data for the query.
([example](https://maps.googleapis.com/maps/api/geocode/json?&address=56%20Exchange%20Terr&components=country:US|administrative_area:Rhode%20Island))

## TODO

- [ ] add optional `direction` query param to the route stops endpoint (to only display trips with the desired direction)
- [ ] add optional `direction` query param to the stop locations endpoint (to limit stops to the desired direction)
- [ ] add additional timing data from the google data matrix api
