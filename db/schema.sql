-- Enable postgis for location queries
CREATE EXTENSION IF NOT EXISTS postgis;

-- Create routes table
CREATE TABLE IF NOT EXISTS routes (
  route_id           integer PRIMARY KEY,
  route_short_name   varchar(5),
  route_long_name    varchar(255),
  route_desc         varchar(255),
  route_type         varchar(255),
  route_url          varchar(255)
);

-- Create trips table
CREATE TABLE IF NOT EXISTS trips (
  trip_id            integer PRIMARY KEY,
  route_id           integer REFERENCES routes(route_id),
  service_id         varchar(255),
  trip_headsign      varchar(255),
  direction_id       bit(1),
  block_id           varchar(255),
  shape_id           varchar(255)
);

-- Create stops table
CREATE TABLE IF NOT EXISTS stops (
  stop_id            varchar(255) PRIMARY KEY,
  stop_code          varchar(255),
  stop_name          varchar(255),
  stop_desc          varchar(255),
  stop_lat           float(8),
  stop_lon           float(8),
  zone_id            varchar(255),
  stop_url           varchar(255),
  location_type      varchar(255),
  parent_station     varchar(255),
  stop_geog          geography(POINT, 4326)
);

-- Build function for creating geography from stop lat and lon
CREATE OR REPLACE FUNCTION build_geog() RETURNS trigger AS $build_geog_trigger$
  BEGIN
    NEW.stop_geog := ST_SetSRID(ST_MakePoint(NEW.stop_lon, NEW.stop_lat), 4326);
    RETURN NEW;
  END;
$build_geog_trigger$ LANGUAGE plpgsql;

-- Add trigger to stops for processing the build geog function
DROP TRIGGER IF EXISTS build_geog_trigger ON stops;
CREATE TRIGGER build_geog_trigger BEFORE INSERT ON stops FOR EACH ROW EXECUTE PROCEDURE build_geog();

-- Create stop times table
CREATE TABLE IF NOT EXISTS stop_times (
trip_id            integer REFERENCES trips(trip_id),
stop_id            varchar(255) REFERENCES stops(stop_id),
arrival_time       varchar(255),
departure_time     varchar(255),
stop_sequence      integer,
pickup_type        varchar(255),
drop_off_type      varchar(255)
);
