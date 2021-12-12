curl -i -XPOST -d '{"id": "123", "attributes":{"a":"a"}, "timestamp": 1}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"id": "123", "attributes":{"a":"a"}, "timestamp": 1}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"id": "123", "attributes":{"a":"a"}, "timestamp": 4}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"id": "123", "attributes":{"a":"a"}, "timestamp": 3}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"id": "123", "attributes":{"b":"b"}, "timestamp": 5}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"id": "123", "attributes":{"b":"b"}, "timestamp": 7}' localhost:4479/api/v1/event
