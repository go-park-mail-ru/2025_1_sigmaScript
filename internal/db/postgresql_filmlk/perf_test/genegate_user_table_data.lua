wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

function generate_unique_hash(thread_id)
  local letters = "abcdefghjklmnopqrstuvwxyz"
  local size = math.random(6, 20)
  local hash = thread_id .. os.time()

  for i = 1, size do
    local random_index = math.random(1, #letters)
    local random_char = letters:sub(random_index, random_index)
    hash = hash .. random_char
  end

  return hash
end

threads = {}

function setup(thread)
  thread:set("thread_id", #threads + 1)
  table.insert(threads, thread)
end

request = function()
  local unique_hash = generate_unique_hash(thread_id)
  local username = "user_" .. unique_hash
  local password_hash = "hashed_password_" .. unique_hash

  local body = string.format(
    '{"username":"%s","password":"%s","repeated_password":"%s"}',
    username,
    password_hash,
    password_hash
  )

  return wrk.format(nil, nil, nil, body)
end
