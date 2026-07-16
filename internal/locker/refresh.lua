local subject = KEYS[1]
local owner = ARGV[1]
local duration = ARGV[2]

local current = redis.call("GET", subject)

if current == false then
	-- key doesn't exist at all
	return -1
elseif current == owner then
	-- owner matches, refresh expiry
	return redis.call("PEXPIRE", subject, duration)
else
	-- key exists but held by a different owner
	return -2
end
