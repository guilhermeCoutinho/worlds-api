-- KEYS[1] = userId
-- ARGV[1] = newWorldId

local userKey = "user:" .. KEYS[1] .. ":world"
local oldWorld = redis.call("GET", userKey)

if oldWorld and oldWorld ~= "" then
    redis.call("SREM", "world:" .. oldWorld .. ":users", KEYS[1])
end

redis.call("SET", userKey, ARGV[1])

local worldKey = "world:" .. ARGV[1] .. ":users"
redis.call("SADD", worldKey, KEYS[1])

return worldKey