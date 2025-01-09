-- 这段代码是一个Redis Lua脚本,用于验证用户输入的验证码
-- 主要逻辑如下:

local key = KEYS[1]                          -- 获取验证码的key
local expectedCode = ARGV[1]                 -- 获取用户输入的验证码
local code = redis.call("get", key)          -- 从Redis获取存储的验证码
local cntKey = key..":cnt"                   -- 验证次数的key
local cnt = tonumber(redis.call("get", cntKey)) -- 获取剩余验证次数并转为数字

if cnt <= 0 then
    -- 如果验证次数小于等于0,说明:
    -- 1. 用户多次输入错误验证码
    -- 2. 验证码已被使用
    -- 返回-1表示验证失败
    return -1
elseif expectedCode == code then
    -- 如果用户输入的验证码正确:
    -- 1. 将验证次数设为-1,防止重复使用
    -- 2. 返回0表示验证成功
    redis.call("set", cntKey, -1)
    return 0
else
    -- 如果验证码错误:
    -- 1. 剩余验证次数减1
    -- 2. 返回-2表示验证码错误
    redis.call("decr", cntKey)
    return -2
end


