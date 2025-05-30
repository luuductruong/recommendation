dto:
	docker compose -f docker-compose.yml up recommendation_proto_builder
proto:
	make -C services/core/proto
