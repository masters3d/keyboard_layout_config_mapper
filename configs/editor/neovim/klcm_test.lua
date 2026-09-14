-- nvim --headless -u NONE -i NONE -n -c "luafile configs/editor/neovim/klcm_test.lua"
package.path = "configs/editor/neovim/?.lua;" .. package.path
local klcm = require("klcm")
local api = vim.api
local notices, checks = {}, 0
vim.notify = function(message) notices[#notices + 1] = message end
vim.o.swapfile = false
vim.o.hidden = true

local function equal(expected, actual, label)
  checks = checks + 1
  if vim.env.KLCM_TEST_TRACE then
    io.stderr:write(checks .. " " .. label .. "\n")
  end
  assert(vim.deep_equal(expected, actual),
    label .. ": expected " .. vim.inspect(expected) .. ", got " .. vim.inspect(actual))
end

local thread
local function resume()
  local ok, err = coroutine.resume(thread)
  if not ok then
    io.stderr:write(tostring(err) .. "\n" .. debug.traceback(thread) .. "\n")
    vim.cmd("cquit 1")
  end
end

local function input(text)
  api.nvim_input(text)
  vim.defer_fn(resume, 20)
  coroutine.yield()
end

local function lines()
  return api.nvim_buf_get_lines(0, 0, -1, false)
end

local function fresh(text, column)
  input("<Esc>")
  vim.cmd("enew!")
  api.nvim_buf_set_lines(0, 0, -1, false, type(text) == "table" and text or { text })
  api.nvim_win_set_cursor(0, { 1, column or 0 })
  input("i")
  equal("i", api.nvim_get_mode().mode, "fresh insert mode")
end

local function signal(key)
  input(key)
end

thread = coroutine.create(function()
  vim.keymap.set("n", "<C-A-F1>", ":echo 'old mapping'<CR>", { silent = true })
  local old = vim.fn.maparg("<C-A-F1>", "n", false, true)
  klcm.setup()
  equal(36, #klcm.actions, "36 semantic actions")
  klcm.setup()
  equal("KLCM move.left", vim.fn.maparg("<C-A-F1>", "n", false, true).desc, "idempotence")

  fresh("abc")
  signal("<C-A-F4>")
  equal({ 1, 1 }, api.nvim_win_get_cursor(0), "insert right")
  input("X")
  equal({ "aXbc" }, lines(), "type at exact insertion boundary")
  signal("<C-A-F8>")
  input("!")
  equal({ "aXbc!" }, lines(), "line end insertion")
  signal("<C-A-F7>")
  input("^")
  equal({ "^aXbc!" }, lines(), "line start insertion")

  fresh("α😀z")
  signal("<C-A-F4>")
  equal({ 1, 2 }, api.nvim_win_get_cursor(0), "Unicode right")
  signal("<C-A-S-F4>")
  equal("s", api.nvim_get_mode().mode, "explicit Select mode")
  input("Q")
  equal({ "αQz" }, lines(), "Unicode selection replaced by typing")
  equal("i", api.nvim_get_mode().mode, "typing leaves Select for Insert")
  signal("<C-A-F1>")
  signal("<C-A-S-F1>")
  input("R")
  equal({ "RQz" }, lines(), "backward Unicode selection")

  fresh("abcd")
  signal("<C-A-S-F4>")
  signal("<C-A-S-F4>")
  input("X")
  equal({ "Xcd" }, lines(), "extend selection")
  fresh("abcd")
  signal("<C-A-S-F4>")
  signal("<C-A-S-F1>")
  equal("i", api.nvim_get_mode().mode, "collapse selection")
  input("X")
  equal({ "Xabcd" }, lines(), "collapsed selection insertion")

  fresh("")
  signal("<C-A-F1><C-A-F4><C-A-F6><C-S-F3>")
  equal({ "" }, lines(), "empty buffer boundaries")
  input("é")
  signal("<C-A-F10>")
  signal("<C-A-F4>")
  input("!")
  equal({ "é!" }, lines(), "EOF Unicode insertion")

  fresh("foo.bar  baz qux")
  signal("<C-A-F6>")
  equal({ 1, 3 }, api.nvim_win_get_cursor(0), "small word punctuation start")
  signal("<C-A-F6>")
  equal({ 1, 4 }, api.nvim_win_get_cursor(0), "small word after punctuation")
  signal("<C-A-F5>")
  equal({ 1, 3 }, api.nvim_win_get_cursor(0), "previous punctuation word")
  vim.fn.setreg('"', "keep unnamed")
  vim.fn.setreg("a", "keep named")
  signal("<C-S-F3>")
  equal({ "foobar  baz qux" }, lines(), "delete current boundary through word start")
  equal("keep unnamed", vim.fn.getreg('"'), "unnamed register preserved")
  equal("keep named", vim.fn.getreg("a"), "named register preserved")

  fresh("one two three four five")
  signal("<C-S-F4>")
  equal({ "four five" }, lines(), "delete three words")
  signal("<C-S-F6>")
  equal({ "" }, lines(), "repeat destructive action")
  signal("<C-S-F1>")
  equal({ "four five" }, lines(), "repeat separate undo unit")
  signal("<C-S-F1>")
  equal({ "one two three four five" }, lines(), "first deletion separate undo unit")
  signal("<C-S-F2>")
  equal({ "four five" }, lines(), "redo")

  fresh("hello world", 2)
  signal("<C-S-F3>")
  equal({ "heworld" }, lines(), "midword delete starts at insertion boundary")
  fresh({ "één", "два three" })
  signal("<C-S-F3>")
  equal({ "два three" }, lines(), "Unicode cross-line delete")

  fresh("before (outer (inner) tail) after", 16)
  signal("<C-S-F5>")
  equal({ "before (outer () tail) after" }, lines(), "innermost parens preserve delimiters")
  input("new")
  equal({ "before (outer (new) tail) after" }, lines(), "parens insertion position")
  fresh("no parentheses", 3)
  signal("<C-S-F5>")
  equal({ "no parentheses" }, lines(), "missing parens safe")
  assert(notices[#notices]:match("No enclosing"))
  fresh("(unmatched", 4)
  signal("<C-S-F5>")
  equal({ "(unmatched" }, lines(), "unmatched parens safe")

  fresh({ "α😀", "xyz", "", "last" })
  signal("<C-A-F4><C-A-F4><C-A-F2>")
  equal({ 2, 2 }, api.nvim_win_get_cursor(0), "vertical Unicode columns")
  signal("<C-A-F12>")
  equal({ 3, 0 }, api.nvim_win_get_cursor(0), "next paragraph")
  signal("<C-A-F10>")
  equal({ 4, 4 }, api.nvim_win_get_cursor(0), "document end")
  signal("<C-A-F11>")
  equal({ 3, 0 }, api.nvim_win_get_cursor(0), "previous paragraph")
  signal("<C-A-F9>")
  equal({ 1, 0 }, api.nvim_win_get_cursor(0), "document start")

  fresh("first second first")
  signal("<C-S-F7>")
  equal("c", api.nvim_get_mode().mode, "visible search command line")
  input("second<Esc>")
  equal({ "first second first" }, lines(), "cancel search does not edit")
  equal("i", api.nvim_get_mode().mode, "cancel search returns insert")
  signal("<C-S-F7>")
  input("first<CR>")
  equal({ 1, 13 }, api.nvim_win_get_cursor(0), "native search accepted")
  signal("<C-S-F8>")
  equal({ 1, 0 }, api.nvim_win_get_cursor(0), "next search wraps")
  signal("<C-S-F9>")
  equal({ 1, 13 }, api.nvim_win_get_cursor(0), "previous search wraps")
  signal("<C-S-F10>")
  assert(notices[#notices]:match("No LSP"))
  signal("<C-S-F12>")
  equal("n", api.nvim_get_mode().mode, "cancel deliberately leaves Normal")
  input("i<Esc>")
  equal("n", api.nvim_get_mode().mode, "Escape stays Normal")
  equal(false, klcm.dispatch("not.an.action"), "unknown action rejected")

  vim.bo.modifiable = false
  equal(false, klcm.dispatch("edit.delete_word"), "unmodifiable excluded")
  vim.bo.modifiable = true
  vim.bo.buftype = "nofile"
  equal(false, klcm.dispatch("move.left"), "special buffer excluded")
  vim.bo.buftype = ""
  klcm.teardown()
  equal(old.rhs, vim.fn.maparg("<C-A-F1>", "n", false, true).rhs, "previous mapping restored")
  equal("", vim.fn.maparg("<C-A-F2>", "i"), "new mapping removed")
  klcm.teardown()
  equal(false, klcm.dispatch("move.left"), "dispatch after teardown rejected")
  klcm.setup({ mappings = { ["move.right"] = "<F6>" }, insert_on_enter = true })
  input(":enew!<CR>")
  equal("i", api.nvim_get_mode().mode, "opt-in buffer entry")
  input("<Esc>")
  equal("n", api.nvim_get_mode().mode, "entry policy does not force InsertEnter")
  input("iabc<F6>X")
  equal({ "abcX" }, lines(), "custom transport")
  input("<Esc>")
  klcm.teardown()
  print("KLCM Neovim: " .. checks .. " assertions passed")
  vim.cmd("qa!")
end)
vim.schedule(resume)
vim.defer_fn(function()
  io.stderr:write("Timed out after " .. checks .. " assertions\n")
  vim.cmd("cquit 1")
end, 15000)
