local concat = table.concat
local kong = kong
local ngx = ngx
local find = string.find
local lower = string.lower
local sha = require "sha2"
local cjson = require("cjson")
local redis = require "resty.redis"
local string = string

local RedisTokenHandler = {}

local sock_opts = {}

local function is_present(str)
    return str and str ~= "" and str ~= null
end

function RedisTokenHandler:access(conf)

    local red = redis:new()
    red:set_timeout(conf.redis_timeout)

    sock_opts.pool = conf.redis_database and
            conf.redis_host .. ":" .. conf.redis_port ..
                    ":" .. conf.redis_database
    local ok, err = red:connect(conf.redis_host, conf.redis_port,
            sock_opts)
    if not ok then
        kong.log.err("failed to connect to Redis: ", err)
        return nil, err
    end

    local times, err = red:get_reused_times()
    if err then
        kong.log.err("failed to get connect reused times: ", err)
        return nil, err
    end

    if times == 0 then
        if is_present(conf.redis_password) then
            local ok, err = red:auth(conf.redis_password)
            if not ok then
                kong.log.err("failed to auth Redis: ", err)
                return nil, err
            end
        end

        if conf.redis_database ~= 0 then
            -- Only call select first time, since we know the connection is shared
            -- between instances that use the same redis database

            local ok, err = red:select(conf.redis_database)
            if not ok then
                kong.log.err("failed to change Redis database: ", err)
                return nil, err
            end
        end
    end

    local token = kong.request.get_header("SUNMI-Token")

    if token == "" then
        return kong.response.exit(500, { message = "token" })
    end

    local current_metric = red:get("tob_admin_Cipher_Turn_" .. token)

    if current_metric == ngx.null then
        return kong.response.exit(500, { message = "token" })
    end

    local id = unserialize(current_metric)
    if id == "" then
        return kong.response.exit(500, { message = current_metric })
    end
    return kong.response.exit(200, id)
end

RedisTokenHandler.PRIORITY = 800
RedisTokenHandler.VERSION = "1.0.0"

function _read_until(data, offset, stopchar)
    --[[
    Read from data[offset] until you encounter some char 'stopchar'.
    ]]
    local buf = {}
    local char = string.sub(data, offset + 1, offset + 1)
    local i = 2
    while not (char == stopchar) do
        -- Consumed all the characters and havent found ';'
        if i + offset > string.len(data) then
            error('Invalid')
        end
        table.insert(buf, char)
        char = string.sub(data, offset + i, offset + i)
        i = i + 1
    end
    -- (chars_read, data)
    return i - 2, table.concat(buf)
end

function _read_chars(data, offset, length)
    --[[
    Read 'length' number of chars from data[offset].
    ]]
    local buf = {}, char
    -- Account for the starting quote char
    -- offset += 1
    for i = 0, length - 1 do
        char = string.sub(data, offset + i, offset + i)
        table.insert(buf, char)
    end
    -- (chars_read, data)
    return length, table.concat(buf)
end

function unserialize(data, offset)
    offset = offset or 0
    --[[
    Find the next token and unserialize it.
    Recurse on array.
    offset = raw offset from start of data
    --]]
    local buf, dtype, dataoffset, typeconvert, datalength, chars, readdata, i,
    key, value, keys, properties, otchars, otype, property
    buf = {}
    dtype = string.lower(string.sub(data, offset + 1, offset + 1))
    -- 't:' = 2 chars
    dataoffset = offset + 2
    typeconvert = function(x)
        return x
    end
    datalength = 0
    chars = datalength
    -- int or double => Number
    if dtype == 'i' or dtype == 'd' then
        typeconvert = function(x)
            return tonumber(x)
        end
        chars, readdata = _read_until(data, dataoffset, ';')
        -- +1 for end semicolon
        dataoffset = dataoffset + chars + 1
        -- bool => Boolean
    elseif dtype == 'b' then
        typeconvert = function(x)
            return tonumber(x) == 1
        end
        chars, readdata = _read_until(data, dataoffset, ';')
        -- +1 for end semicolon
        dataoffset = dataoffset + chars + 1
        -- n => None
    elseif dtype == 'n' then
        readdata = nil
        -- s => String
    elseif dtype == 's' then
        chars, stringlength = _read_until(data, dataoffset, ':')
        -- +2 for colons around length field
        dataoffset = dataoffset + chars + 2
        -- +1 for start quote
        chars, readdata = _read_chars(data, dataoffset + 1, tonumber(stringlength))
        -- +2 for endquote semicolon
        dataoffset = dataoffset + chars + 2
        --[[
        TODO
        review original: if chars != int(stringlength) != int(readdata):
        ]]
        if not (chars == tonumber(stringlength)) then
            error('String length mismatch')
        end
        -- array => Table
        -- If you originally serialized a Tuple or List, it will
        -- be unserialized as a Dict.  PHP doesn't have tuples or lists,
        -- only arrays - so everything has to get converted into an array
        -- when serializing and the original type of the array is lost
    elseif dtype == 'a' then
        readdata = {}
        -- How many keys does this list have?
        chars, keys = _read_until(data, dataoffset, ':')
        -- +2 for colons around length field
        dataoffset = dataoffset + chars + 2
        -- Loop through and fetch this number of key/value pairs
        for i = 0, tonumber(keys) - 1 do
            -- Read the key
            key, ktype, kchars = unserialize(data, dataoffset)
            dataoffset = dataoffset + kchars
            -- Read value of the key
            value, vtype, vchars = unserialize(data, dataoffset)
            -- Cound ending bracket of nested array
            if vtype == 'a' then
                vchars = vchars + 1
            end
            dataoffset = dataoffset + vchars
            -- Set the list element
            readdata[key] = value
        end
        -- object => Table
    elseif dtype == 'o' then
        readdata = {}
        -- How log is the type of this object?
        chars, otchars = _read_until(data, dataoffset, ':')
        dataoffset = dataoffset + chars + 2
        -- Which type is this object?
        otype = string.sub(data, dataoffset + 1, dataoffset + otchars)
        dataoffset = dataoffset + otchars + 2
        if otype == 'stdClass' then
            -- How many properties does this list have?
            chars, properties = _read_until(data, dataoffset, ':')
            -- +2 for colons around length field
            dataoffset = dataoffset + chars + 2
            -- Loop through and fetch this number of key/value pairs
            for i = 0, tonumber(properties) - 1 do
                -- Read the key
                property, ktype, kchars = unserialize(data, dataoffset)
                dataoffset = dataoffset + kchars
                -- Read value of the key
                value, vtype, vchars = unserialize(data, dataoffset)
                -- Cound ending bracket of nested array
                if vtype == 'a' then
                    vchars = vchars + 1
                end
                dataoffset = dataoffset + vchars
                -- Set the list element
                readdata[property] = value
            end
        else
            _unknown_type(dtype)
        end
    else
        _unknown_type(dtype)
    end
    --~ return (dtype, dataoffset-offset, typeconvert(readdata))
    return typeconvert(readdata), dtype, dataoffset - offset
end
-- I don't know how to unserialize this



return RedisTokenHandler
