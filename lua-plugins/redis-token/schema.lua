local typedefs = require "kong.db.schema.typedefs"

return {
    name = "redis-token",
    fields = {
        { consumer = typedefs.no_consumer },
        { protocols = typedefs.protocols_http },
        { config = {
            type = "record",
            fields = {
                { redis_host = typedefs.host },
                { redis_port = typedefs.port({ default = 6379 }), },
                { redis_password = { type = "string", len_min = 0 }, },
                { redis_timeout = { type = "number", default = 2000, }, },
                { redis_database = { type = "integer", default = 0 }, },
                { php_serialization = { type = "boolean", default = false }, },
                { token_name = { type = "string", default = "X-Token" }, },
                { id_name = { type = "string", default = "X-Id" }, },
            },
        },
        },
    },
}