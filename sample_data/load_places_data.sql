drop table if exists place;

create table place (
    "food_resource_type" text,
    "agency" text,
    "location" text,
    "operational_status" text,
    "operational_notes" text,
    "who_they_serve" text,
    "address" text,
    "latitude" text,
    "longitude" text,
    "phone_number" text,
    "website" text,
    "days_or_hours" text,
    "date_updated" text,
    "tsv"         tsvector
);

create trigger tsvectorupdate before insert or update
on place for each row execute procedure
tsvector_update_trigger(
	tsv, 'pg_catalog.english', "agency", "location", "operational_notes"
);

create index index_pages_on_tsv on place using gin (tsv);

\copy place("food_resource_type", "agency", "location", "operational_status", "operational_notes", "who_they_serve", "address", "latitude", "longitude", "phone_number", "website", "days_or_hours", "date_updated") from '/tmp/sample_data/Emergency_Food_and_Meals_Seattle_and_King_County.csv' with null as E'\'\'' delimiter ',' CSV HEADER
