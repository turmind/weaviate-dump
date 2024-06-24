# weaviate-dump

Export Weaviate data to Excel by Weaviate Restful API

## flags

params|required|type|default value|note
-|-|-|-|-
host|YES|String|localhost|weaviate host
port|YES|Integer|8080|weaviate port
token|NO|String|null|weaviate access token
class|YES|String|null|the class that needs to be dumped
limit|YES|Integer|25|The size of each request

## demo

```shell
weaviate-dump --host 127.0.0.1 --port 8080 --token token_xxxxxxxxxxxx --class class_xxxxxxxxxxxx --limit 10
```