-- Opt in with require("klcm").setup({ insert_on_enter = false }).
-- opts.mappings replaces defaults: { ["move.left"] = "<C-A-F1>", ... }.
-- Positions are UTF-8 codepoint boundaries, not grapheme/display columns.
local M = {}
local api, fn = vim.api, vim.fn
local state
local targets = {
  "left", "down", "up", "right", "word_previous", "word_next",
  "line_start", "line_end", "document_start", "document_end",
  "paragraph_previous", "paragraph_next",
}
local edits = {
  "edit.undo", "edit.redo", "edit.delete_word", "edit.delete_3_words",
  "edit.change_inside_parens", "edit.repeat", "search.open", "search.next",
  "search.previous", "code.definition", "history.back", "interaction.cancel",
}

local function keys(text)
  return api.nvim_replace_termcodes(text, true, false, true)
end

local function notice(text)
  vim.notify("KLCM: " .. text, vim.log.levels.WARN)
  return false
end

local function editable()
  return vim.bo.buftype == "" and vim.bo.modifiable and not vim.bo.readonly
end

local function document()
  local lines = api.nvim_buf_get_lines(0, 0, -1, false)
  return lines, table.concat(lines, "\n")
end

local function offset(lines, cursor)
  local n = 0
  for row = 1, cursor[1] - 1 do
    n = n + #lines[row] + 1
  end
  return n + math.min(cursor[2], #lines[cursor[1]])
end

local function position(lines, n)
  for row, line in ipairs(lines) do
    if n <= #line or row == #lines then
      return { row, math.min(n, #line) }
    end
    n = n - #line - 1
  end
end

local function next_char(text, n)
  n = math.min(n + 1, #text)
  while n < #text and text:byte(n + 1) >= 128 and text:byte(n + 1) < 192 do
    n = n + 1
  end
  return n
end

local function previous_char(text, n)
  n = math.max(n - 1, 0)
  while n > 0 and text:byte(n + 1) >= 128 and text:byte(n + 1) < 192 do
    n = n - 1
  end
  return n
end

local function class(text, n)
  if n >= #text then
    return "end"
  end
  local c = text:sub(n + 1, next_char(text, n))
  if c:match("%s") then
    return "space"
  end
  return fn.match(c, [[\k]]) >= 0 and "keyword" or "punctuation"
end

local function word_next(text, n)
  local kind = class(text, n)
  if kind ~= "space" then
    while n < #text and class(text, n) == kind do
      n = next_char(text, n)
    end
  end
  while n < #text and class(text, n) == "space" do
    -- Like w, stop on an empty line rather than skipping all paragraphs.
    if text:sub(n + 1, n + 2) == "\n\n" then
      return n + 1
    end
    n = next_char(text, n)
  end
  return n
end

local function word_previous(text, n)
  if n == 0 then
    return n
  end
  n = previous_char(text, n)
  while n > 0 and class(text, n) == "space" do
    if text:sub(n, n + 1) == "\n\n" then
      return n
    end
    n = previous_char(text, n)
  end
  local kind = class(text, n)
  while n > 0 and class(text, previous_char(text, n)) == kind do
    n = previous_char(text, n)
  end
  return n
end

local function move(target, lines, text, n)
  local p = position(lines, n)
  if target == "left" then
    return previous_char(text, n)
  elseif target == "right" then
    return next_char(text, n)
  elseif target == "word_next" then
    return word_next(text, n)
  elseif target == "word_previous" then
    return word_previous(text, n)
  elseif target == "line_start" then
    return n - p[2]
  elseif target == "line_end" then
    return n - p[2] + #lines[p[1]]
  elseif target == "document_start" then
    return 0
  elseif target == "document_end" then
    return #text
  elseif target == "up" or target == "down" then
    local row = math.max(1, math.min(#lines, p[1] + (target == "up" and -1 or 1)))
    local column = fn.strchars(lines[p[1]]:sub(1, p[2]))
    return offset(lines, { row, #fn.strcharpart(lines[row], 0, column) })
  elseif target == "paragraph_previous" or target == "paragraph_next" then
    local step = target == "paragraph_previous" and -1 or 1
    local row = p[1] + step
    while row > 1 and row < #lines and lines[row] ~= "" do
      row = row + step
    end
    if row < 1 then
      return 0
    elseif row > #lines or (row == #lines and lines[row] ~= "") then
      return #text
    end
    return offset(lines, { row, 0 })
  end
end

local function clear_selection()
  if state.selection then
    vim.o.selection = state.selection.option
    state.selection = nil
  end
end

local function insert_at(lines, n)
  local p = position(lines, n)
  api.nvim_win_set_cursor(0, p)
  if p[2] == #lines[p[1]] then
    vim.cmd("startinsert!")
  else
    vim.cmd("startinsert")
  end
end

local function select_range(lines, text, anchor, active)
  if anchor == active then
    insert_at(lines, active)
    return
  end
  local selection = {
    anchor = anchor, active = active, buffer = api.nvim_get_current_buf(),
    option = vim.o.selection,
  }
  vim.o.selection = "inclusive"
  local low, high = math.min(anchor, active), math.max(anchor, active)
  local first, last = low, previous_char(text, high)
  if active < anchor then
    first, last = last, first
  end
  api.nvim_win_set_cursor(0, position(lines, first))
  vim.cmd("normal! v")
  api.nvim_win_set_cursor(0, position(lines, last))
  vim.cmd("normal! " .. keys("<C-g>"))
  state.selection = selection
end

local function delete_range(lines, first, last)
  if first == last then
    return false
  end
  local a, b = position(lines, first), position(lines, last)
  -- Each API edit starts a fresh undo block; never use :undojoin.
  api.nvim_buf_set_text(0, a[1] - 1, a[2], b[1] - 1, b[2], {})
  return true
end

local function parens(text, n)
  -- Literal balanced parentheses only; strings/comments are not parsed.
  local stack, best = {}, nil
  for i = 1, #text do
    local c = text:sub(i, i)
    if c == "(" then
      stack[#stack + 1] = i - 1
    elseif c == ")" and #stack > 0 then
      local open = table.remove(stack)
      if open <= n and n <= i - 1 and (not best or open > best[1]) then
        best = { open + 1, i - 1 }
      end
    end
  end
  return best
end

local function search(direction)
  local pattern = fn.getreg("/")
  if pattern == "" then
    return notice("No search pattern; use search.open first")
  end
  local ok, found = pcall(fn.search, pattern, direction == "previous" and "b" or "")
  if not ok or found == 0 then
    return notice(ok and "Search pattern not found" or "Invalid search pattern")
  end
  return true
end

local actions = {}
for _, target in ipairs(targets) do
  actions["move." .. target] = true
  actions["select." .. target] = true
end
for _, id in ipairs(edits) do
  actions[id] = true
end
M.actions = vim.tbl_keys(actions)
table.sort(M.actions)

local function run(id, context)
  if not state or not editable() then
    return notice("Action unavailable in this buffer")
  end
  local lines, text = document()
  local n = math.min(context.offset, #text)
  local target = id:match("^move%.(.+)") or id:match("^select%.(.+)")
  if target then
    local dest = move(target, lines, text, n)
    if id:match("^select%.") then
      select_range(lines, text, context.anchor or n, dest)
    else
      insert_at(lines, dest)
    end
    return true
  end
  if id == "interaction.cancel" then
    api.nvim_win_set_cursor(0, position(lines, n))
    return true
  elseif id == "search.open" then
    state.prompt = true
    -- Native / command line is visible and cancellable; no Ex interpolation.
    api.nvim_feedkeys("/", "n", false)
    return true
  elseif id == "search.next" or id == "search.previous" then
    api.nvim_win_set_cursor(0, position(lines, n))
    local ok = search(id:match("%.([^.]+)$"))
    insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
    return ok
  elseif id == "edit.undo" or id == "edit.redo" then
    vim.cmd(id == "edit.undo" and "silent undo" or "silent redo")
    lines = document()
    insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
    return true
  elseif id == "code.definition" then
    local clients = vim.lsp.get_clients and vim.lsp.get_clients({ bufnr = 0 })
      or vim.lsp.get_active_clients({ bufnr = 0 })
    for _, client in ipairs(clients) do
      if client.supports_method("textDocument/definition") then
        vim.lsp.buf.definition()
        return true
      end
    end
    insert_at(lines, n)
    return notice("No LSP definition provider attached")
  elseif id == "history.back" then
    local jumps = fn.getjumplist()
    if jumps[2] == 0 then
      insert_at(lines, n)
      return notice("No earlier jump")
    end
    vim.cmd("normal! " .. keys("<C-o>"))
    if editable() then
      lines = document()
      insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
    end
    return true
  end
  if id == "edit.repeat" then
    if not state.last_edit then
      insert_at(lines, n)
      return notice("No completed destructive action to repeat")
    end
    id = state.last_edit
  end
  local first, last = n, n
  if id == "edit.delete_word" or id == "edit.delete_3_words" then
    for _ = 1, id == "edit.delete_word" and 1 or 3 do
      last = word_next(text, last)
    end
  elseif id == "edit.change_inside_parens" then
    local range = parens(text, n)
    if not range then
      insert_at(lines, n)
      return notice("No enclosing balanced parentheses")
    end
    first, last = range[1], range[2]
  end
  if delete_range(lines, first, last) then
    state.last_edit = id
  end
  lines = document()
  insert_at(lines, first)
  return true
end

function M.dispatch(id)
  if vim.env.KLCM_TEST_TRACE then io.stderr:write("dispatch " .. id .. "\n") end
  if not state then
    return notice("Call setup() before dispatching actions")
  end
  if not actions[id] then
    return notice("Unsupported action: " .. tostring(id))
  end
  if not editable() then
    return notice("Action unavailable in this buffer")
  end
  local mode = api.nvim_get_mode().mode
  local lines = document()
  local context = { offset = offset(lines, api.nvim_win_get_cursor(0)) }
  if mode:match("^[sS]") and state.selection
      and state.selection.buffer == api.nvim_get_current_buf() then
    context.offset = state.selection.active
    context.anchor = state.selection.anchor
  end
  clear_selection()
  if mode:match("^[isS]") then
    -- One explicit exit before the action, not a queued editing macro.
    state.pending = { id = id, context = context, buffer = api.nvim_get_current_buf() }
    api.nvim_feedkeys(keys("<C-\\><C-n><Cmd>lua require('klcm')._resume()<CR>"), "ni", false)
    if vim.env.KLCM_TEST_TRACE then io.stderr:write("queued\n") end
    return true
  end
  return run(id, context)
end

function M._resume()
  if vim.env.KLCM_TEST_TRACE then io.stderr:write("resume\n") end
  if not state or not state.pending then
    return
  end
  local pending = state.pending
  state.pending = nil
  if pending.buffer == api.nvim_get_current_buf() then
    run(pending.id, pending.context)
  end
end

function M.teardown()
  if not state then
    return
  end
  clear_selection()
  for _, map in ipairs(state.maps) do
    pcall(vim.keymap.del, map.mode, map.lhs)
    if map.previous and next(map.previous) then
      fn.mapset(map.mode, false, map.previous)
    end
  end
  api.nvim_del_augroup_by_id(state.group)
  state = nil
end

function M.setup(opts)
  M.teardown()
  opts = opts or {}
  state = { maps = {}, group = api.nvim_create_augroup("KLCMReference", { clear = true }) }
  local mappings = opts.mappings
  if not mappings then
    mappings = {}
    for i, target in ipairs(targets) do
      mappings["move." .. target] = "<C-A-F" .. i .. ">"
      mappings["select." .. target] = "<C-A-S-F" .. i .. ">"
      mappings[edits[i]] = "<C-S-F" .. i .. ">"
    end
  end
  for id, lhs in pairs(mappings) do
    if not actions[id] or type(lhs) ~= "string" then
      notice("Invalid transport mapping: " .. tostring(id))
    else
      for _, mode in ipairs({ "i", "n", "s" }) do
        local previous = fn.maparg(lhs, mode, false, true)
        -- Buffer-local mappings are left untouched and continue to take precedence.
        if previous.buffer == 1 then
          previous = {}
          for _, map in ipairs(api.nvim_get_keymap(mode)) do
            if keys(map.lhs) == keys(lhs) then
              previous = map
              break
            end
          end
        end
        state.maps[#state.maps + 1] = { mode = mode, lhs = lhs, previous = previous }
        vim.keymap.set(mode, lhs, function() M.dispatch(id) end, { desc = "KLCM " .. id })
      end
    end
  end
  api.nvim_create_autocmd("ModeChanged", {
    group = state.group,
    callback = function()
      if not api.nvim_get_mode().mode:match("^[sS]") then
        clear_selection()
      end
    end,
  })
  api.nvim_create_autocmd("CmdlineLeave", {
    group = state.group,
    pattern = "/",
    callback = function()
      if state.prompt then
        state.prompt = false
        vim.schedule(function()
          if state and editable() and api.nvim_get_mode().mode == "n" then
            local lines = document()
            insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
          end
        end)
      end
    end,
  })
  if opts.insert_on_enter then
    api.nvim_create_autocmd("BufEnter", {
      group = state.group,
      callback = function()
        if editable() and api.nvim_get_mode().mode == "n" then
          vim.cmd("startinsert")
        end
      end,
    })
  end
  return M
end

return M
