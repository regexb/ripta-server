#!/bin/sh

DATABASE_NAME="${DATABASE_NAME:-ripta}"
DATABASE_USER="${DATABASE_USER:-$USER}"
DATABASE_HOST="${DATABASE_HOST:-localhost}"

# load schema
psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -a -f schema.sql

# import csv files from the google transit exports (http://www.ripta.com/googledata/current/google_transit.zip)
psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -c "\copy routes (route_id,route_short_name,route_long_name,route_desc,route_type,route_url) FROM '$(pwd)/google_transit/routes.txt' CSV HEADER;"
psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -c "\copy stops (stop_id,stop_code,stop_name,stop_desc,stop_lat,stop_lon,zone_id,stop_url,location_type,parent_station) FROM '$(pwd)/google_transit/stops.txt' CSV HEADER;"
psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -c "\copy trips (route_id,service_id,trip_id,trip_headsign,direction_id,block_id,shape_id) FROM '$(pwd)/google_transit/trips.txt' CSV HEADER;"
psql -h $DATABASE_HOST -U $DATABASE_USER -d $DATABASE_NAME -c "\copy stop_times (trip_id,arrival_time,departure_time,stop_id,stop_sequence,pickup_type,drop_off_type) FROM '$(pwd)/google_transit/stop_times.txt' CSV HEADER;"
