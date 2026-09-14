-- Opt in with require("klcm").setup({ insert_on_enter = false }).
-- opts.mappings replaces defaults: { ["move.left"] = "<C-A-F1>", ... }.
-- Positions are UTF-8 codepoint boundaries, not grapheme/display columns.
-- dispatch(id) uses the active selection endpoint. Repeat replays the last
-- successful delete/paren-clear, not text subsequently typed or native Vim ".".
-- Search accepts Vim regex, not /offsets or chained native search commands.
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
  local kind = fn.charclass(c)
  if c:match("%s") or kind == 0 then
    return "space"
  end
  return kind
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
    return p[2] == 0 and n or previous_char(text, n)
  elseif target == "right" then
    return p[2] == #lines[p[1]] and n or next_char(text, n)
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
    api.nvim_win_set_cursor(0, p)
    vim.cmd("normal! " .. (target == "paragraph_previous" and "{" or "}"))
    local dest = api.nvim_win_get_cursor(0)
    if target == "paragraph_next" and dest[1] == #lines
        and dest[2] >= previous_char(lines[#lines], #lines[#lines]) then
      return #text
    end
    return offset(lines, dest)
  end
end

local function clear_selection()
  if state.selection then
    vim.o.selection = state.selection.option
    if api.nvim_win_is_valid(state.selection.window) then
      vim.wo[state.selection.window].virtualedit = state.selection.virtualedit
    end
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
    option = vim.o.selection, window = api.nvim_get_current_win(),
    virtualedit = vim.wo.virtualedit,
  }
  vim.o.selection = "inclusive"
  vim.wo.virtualedit = "onemore"
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
  -- Close the undo block even for consecutive programmatic Normal-mode calls.
  vim.bo.undolevels = vim.bo.undolevels
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

local function search(direction, include_current)
  local pattern = fn.getreg("/")
  if pattern == "" then
    return notice("No search pattern; use search.open first")
  end
  local flags = direction == "previous" and "b" or ""
  local ok, found = pcall(fn.search, pattern, flags .. (include_current and "c" or ""))
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

local install_prompt_mapping, restore_prompt_mapping

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
    insert_at(lines, n)
    return true
  elseif id == "search.open" then
    state.prompt = {
      buffer = api.nvim_get_current_buf(), window = api.nvim_get_current_win(), offset = n,
      tick = api.nvim_buf_get_changedtick(0),
    }
    api.nvim_win_set_cursor(0, position(lines, n))
    install_prompt_mapping()
    -- The native / prompt is visible; submitted regex is handled without Ex.
    api.nvim_feedkeys("/", "ni", false)
    return true
  elseif id == "search.next" or id == "search.previous" then
    local p = position(lines, n)
    api.nvim_win_set_cursor(0, p)
    local backwards_from_end = id == "search.previous" and p[2] > 0 and p[2] == #lines[p[1]]
    local ok = search(id:match("%.([^.]+)$"), backwards_from_end)
    insert_at(lines, ok and offset(lines, api.nvim_win_get_cursor(0)) or n)
    return ok
  elseif id == "edit.undo" or id == "edit.redo" then
    vim.cmd(id == "edit.undo" and "silent undo" or "silent redo")
    lines = document()
    insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
    return true
  elseif id == "code.definition" then
    api.nvim_win_set_cursor(0, position(lines, n))
    local clients = vim.lsp.get_clients and vim.lsp.get_clients({ bufnr = 0 })
      or vim.lsp.get_active_clients({ bufnr = 0 })
    for _, client in ipairs(clients) do
      if client.supports_method("textDocument/definition") then
        local ok, err = pcall(vim.lsp.buf.definition)
        if ok then
          if editable() then
            lines = document()
            insert_at(lines, offset(lines, api.nvim_win_get_cursor(0)))
          end
          return true
        end
        insert_at(lines, n)
        return notice("Definition request failed: " .. tostring(err))
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
  if not state then
    return notice("Call setup() before dispatching actions")
  end
  if not actions[id] then
    return notice("Unsupported action: " .. tostring(id))
  end
  if not editable() then
    return notice("Action unavailable in this buffer")
  end
  if state.pending then
    return notice("Previous action is still pending")
  end
  local mode = api.nvim_get_mode().mode
  if mode == "c" and id == "interaction.cancel" and state.prompt then
    api.nvim_feedkeys(keys("<C-c>"), "ni", false)
    return true
  elseif mode ~= "n" and not mode:match("^[isS]") then
    return notice("Action unavailable in mode " .. mode)
  end
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
    state.pending = {
      id = id, context = context, buffer = api.nvim_get_current_buf(),
      window = api.nvim_get_current_win(), tick = api.nvim_buf_get_changedtick(0),
    }
    api.nvim_feedkeys(keys("<C-\\><C-n><Cmd>lua require('klcm')._resume()<CR>"), "ni", false)
    return true
  end
  return run(id, context)
end

function M._resume()
  if not state or not state.pending then
    return
  end
  local pending = state.pending
  state.pending = nil
  if pending.buffer ~= api.nvim_get_current_buf() or pending.window ~= api.nvim_get_current_win()
      or pending.tick ~= api.nvim_buf_get_changedtick(0) then
    notice("Action cancelled: buffer, window, or text changed while leaving Insert/Select mode")
    if editable() then
      local lines, text = document()
      local n = offset(lines, api.nvim_win_get_cursor(0))
      while n > 0 and n < #text and text:byte(n + 1) >= 128 and text:byte(n + 1) < 192 do
        n = n - 1
      end
      insert_at(lines, n)
    end
    return
  end
  run(pending.id, pending.context)
end

local function local_mapping(buffer, mode, lhs)
  for _, map in ipairs(api.nvim_buf_get_keymap(buffer, mode)) do
    if keys(map.lhs) == keys(lhs) then
      return map
    end
  end
end

local function restore_mapping(map)
  if not api.nvim_buf_is_valid(map.buffer) then
    return
  end
  local current = local_mapping(map.buffer, map.mode, map.lhs)
  -- Do not remove a mapping installed by the user after setup().
  if current and current.callback == map.callback then
    vim.keymap.del(map.mode, map.lhs, { buffer = map.buffer })
    if map.previous then
      api.nvim_buf_call(map.buffer, function() fn.mapset(map.mode, false, map.previous) end)
    end
  end
end

local function detach(buffer)
  for _, map in ipairs(state.maps[buffer] or {}) do
    restore_mapping(map)
  end
  state.maps[buffer] = nil
end

local function refresh(buffer)
  if not state or not api.nvim_buf_is_loaded(buffer) then
    return
  end
  local eligible = vim.bo[buffer].buftype == "" and vim.bo[buffer].modifiable and not vim.bo[buffer].readonly
  if not eligible then
    detach(buffer)
  else
    detach(buffer)
    state.maps[buffer] = {}
    for id, lhs in pairs(state.mappings) do
      for _, mode in ipairs({ "i", "n", "s" }) do
        if not local_mapping(buffer, mode, lhs) then
          local callback = function()
            if not editable() then
              -- Also handles API option changes that do not emit OptionSet.
              detach(api.nvim_get_current_buf())
              api.nvim_feedkeys(keys(lhs), "mi", false)
              return
            end
            M.dispatch(id)
          end
          state.maps[buffer][#state.maps[buffer] + 1] = {
            buffer = buffer, mode = mode, lhs = lhs, callback = callback,
          }
          vim.keymap.set(mode, lhs, callback, { buffer = buffer, desc = "KLCM " .. id })
        end
      end
    end
  end
end

restore_prompt_mapping = function()
  if state.prompt_map then
    restore_mapping(state.prompt_map)
    state.prompt_map = nil
  end
end

install_prompt_mapping = function()
  local lhs = state.mappings["interaction.cancel"]
  if not lhs then
    return
  end
  local buffer = api.nvim_get_current_buf()
  local callback = function() return "<C-c>" end
  state.prompt_map = {
    buffer = buffer, mode = "c", lhs = lhs, callback = callback,
    previous = local_mapping(buffer, "c", lhs),
  }
  vim.keymap.set("c", lhs, callback, { buffer = buffer, expr = true, desc = "KLCM interaction.cancel" })
end

function M._finish_search()
  if not state or not state.finished_prompt then
    return
  end
  local prompt = state.finished_prompt
  state.finished_prompt = nil
  if api.nvim_get_current_win() == prompt.window and api.nvim_get_current_buf() == prompt.buffer
      and editable() and api.nvim_get_mode().mode == "n" then
    local lines = document()
    local unchanged = prompt.tick == api.nvim_buf_get_changedtick(0)
    local original = unchanged and prompt.offset or offset(lines, api.nvim_win_get_cursor(0))
    local found = false
    if not prompt.aborted then
      if prompt.pattern ~= "" then
        fn.setreg("/", prompt.pattern)
        fn.histadd("/", prompt.pattern)
      end
      api.nvim_win_set_cursor(0, position(lines, original))
      found = search("next")
    end
    insert_at(lines, found and offset(lines, api.nvim_win_get_cursor(0)) or original)
  end
end

function M.teardown()
  if not state then
    return
  end
  if state.selection then
    vim.cmd("normal! " .. keys("<C-\\><C-n>"))
    clear_selection()
  end
  restore_prompt_mapping()
  for buffer in pairs(state.maps) do
    detach(buffer)
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
  state.mappings = {}
  for id, lhs in pairs(mappings) do
    if not actions[id] or type(lhs) ~= "string" then
      notice("Invalid transport mapping: " .. tostring(id))
    else
      state.mappings[id] = lhs
    end
  end
  for _, buffer in ipairs(api.nvim_list_bufs()) do
    refresh(buffer)
  end
  api.nvim_create_autocmd({ "BufEnter", "BufWinEnter" }, {
    group = state.group, callback = function(event) refresh(event.buf) end,
  })
  api.nvim_create_autocmd("OptionSet", {
    group = state.group, pattern = { "buftype", "modifiable", "readonly" },
    callback = function() refresh(api.nvim_get_current_buf()) end,
  })
  api.nvim_create_autocmd("ModeChanged", {
    group = state.group,
    callback = function()
      if not api.nvim_get_mode().mode:match("^[sS]") then
        clear_selection()
      end
    end,
  })
  api.nvim_create_autocmd("BufLeave", { group = state.group, callback = clear_selection })
  api.nvim_create_autocmd("CmdlineLeave", {
    group = state.group,
    pattern = "/",
    callback = function()
      if state.prompt then
        local prompt = state.prompt
        prompt.aborted = vim.v.event.abort
        prompt.pattern = fn.getcmdline()
        -- Native search errors flush typeahead. Cancel native execution and
        -- evaluate the regex with search() under pcall after the prompt closes.
        vim.cmd("let v:event.abort = v:true")
        state.prompt = nil
        restore_prompt_mapping()
        state.finished_prompt = prompt
        api.nvim_feedkeys(keys("<Cmd>lua require('klcm')._finish_search()<CR>"), "ni", false)
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
