
-- 发送到redis的key
local key = KEYS[1]

-- 发送到redis的value
local cntKey = key .. ":cnt"
local value = ARGV[1]

-- TTL返回值可能为nil,需要判断
local ttl = redis.call("TTL", key)
if ttl == nil then
    return -2
end

ttl = tonumber(ttl)

-- -1是key存在，但没设置过期时间
if ttl == -1 then
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("SET", key, value)
    redis.call("EXPIRE", key, 600)
    redis.call("SET", cntKey, 3)
    redis.call("EXPIRE", cntKey, 600)
    return 0
else
    -- 已经发了一个验证码，不到1分钟
    return -1
end