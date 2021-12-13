curl -i -XPOST -d '{"id": "1", "attributes": {"a":"a"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"a":"a"}, "startTimestamp": 1, "endTimestamp": 59}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"b":"b"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"c":"c"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"a":"a","b":"b"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"c":"c","b":"b"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
curl -i -XPOST -d '{"id": "1", "attributes": {"c":"c2","b":"b"}, "startTimestamp": 1, "endTimestamp": 300}' http://localhost:4479/api/v1/query
