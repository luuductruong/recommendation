dto:
	docker compose -f docker-compose.yml up recommendation_proto_builder
dbs:
	docker compose -f docker-compose.yml up -d recommendation_postgres
proto:
	make -C services/core/proto
