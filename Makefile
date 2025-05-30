dto:
	docker compose -f docker-compose.yml up recommendation_proto_builder recommendation_postgres
proto:
	make -C services/core/proto
