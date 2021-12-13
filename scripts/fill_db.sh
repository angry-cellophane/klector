curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"a":"a"}, "timestamp": 30}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"a":"a"}, "timestamp": 32}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"a":"a"}, "timestamp": 56}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"a":"a"}, "timestamp": 110}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"a":"a", "b":"b"}, "timestamp": 90}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"b":"b"}, "timestamp": 40}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"b":"b"}, "timestamp": 23}]}' localhost:4479/api/v1/event
curl -i -XPOST -d '{"events":[{"id": "123", "attributes":{"b":"b", "c":"c"}, "timestamp": 240}]}' localhost:4479/api/v1/event