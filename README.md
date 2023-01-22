Data Source: http://www.seattle.gov/humanservices/

Read along in blog / guide-ish [https://home.gerardbentley.com/guides/golang_postgres_stack/](https://home.gerardbentley.com/guides/golang_postgres_stack/)

```sh

docker-compose up --build
docker-compose down --volumes --remove-orphans
docker-compose exec database psql -U places_user -d places -f /tmp/sample_data/load_places_data.sql
curl -v "http://127.0.0.1:5000/place?name=First" | jq '.'
curl -v "http://127.0.0.1:5000/place?name=nonsense"
```

```psql
docker-compose exec database psql -U places_user -d places
\dt
select row_to_json(row) from (select * from place limit 1) row;
select count(*) from place where tsv @@ to_tsquery('Community');
```
