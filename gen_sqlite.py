#!/usr/bin/python

import sys
import csv

print """CREATE TABLE routes ( 
            route_id STRING PRIMARY KEY,
            agency_id STRING,
            route_short_name STRING,
            route_long_name STRING,
            route_desc STRING,
            route_type STRING,
            route_url STRING,
            route_color STRING,
            route_text_color STRING
);"""
with open('/'.join([sys.argv[1], 'routes.txt'])) as csv_file:
    reader = csv.reader(csv_file)
    header = reader.next()
    for row in reader:
        print "INSERT INTO routes (" + ", ".join(header) + ") VALUES " + '("' + '", "'.join(row) + '");'

print """CREATE TABLE stops (
            stop_id STRING,
            stop_code STRING,
            stop_name STRING,
            stop_desc STRING,
            stop_lat STRING,
            stop_lon STRING,
            zone_id STRING,
            stop_url STRING,
            location_type STRING,
            parent_station STRING
        );"""
with open('/'.join([sys.argv[1], 'stops.txt'])) as csv_file:
    reader = csv.reader(csv_file)
    header = reader.next()
    for row in reader:
        print "INSERT INTO stops (" + ", ".join(header) + ") VALUES " + '("' + '", "'.join(row) + '");'

print """CREATE TABLE stop_times (
            trip_id STRING,
            arrival_time STRING,
            departure_time STRING,
            stop_id STRING,
            stop_sequence NUMBER,
            stop_headsign NUMBER,
            pickup_type STRING,
            drop_off_type STRING,
            shape_dist_traveled STRING
        );"""
with open('/'.join([sys.argv[1], 'stop_times.txt'])) as csv_file:
    reader = csv.reader(csv_file)
    header = reader.next()
    for row in reader:
        print "INSERT INTO stop_times (" + ", ".join(header) + ") VALUES " + '("' + '", "'.join(row) + '");'

print """CREATE TABLE trips (
            route_id STRING,
            service_id STRING,
            trip_id STRING,
            trip_headsign STRING,
            direction_id STRING,
            block_id STRING,
            shape_id STRING
        );"""
with open('/'.join([sys.argv[1], 'trips.txt'])) as csv_file:
    reader = csv.reader(csv_file)
    header = reader.next()
    for row in reader:
        print "INSERT INTO trips (" + ", ".join(header) + ") VALUES " + '("' + '", "'.join(row) + '");'

print """CREATE TABLE transfers (
            from_stop_id STRING,
            to_stop_id STRING,
            transfer_type STRING,
            min_transfer_time NUMBER
        );"""
with open('/'.join([sys.argv[1], 'transfers.txt'])) as csv_file:
    reader = csv.reader(csv_file)
    header = reader.next()
    for row in reader:
        print "INSERT INTO transfers (" + ", ".join(header) + ") VALUES " + '("' + '", "'.join(row) + '");'
