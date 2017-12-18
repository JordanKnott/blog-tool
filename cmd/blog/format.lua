--
-- Created by IntelliJ IDEA.
-- User: NightWolf
-- Date: 12/18/2017
-- Time: 5:55 AM
-- To change this template use File | Settings | File Templates.
--

json = require("json")

function format(title, date, flags_json)
    local flags = json.decode(flags_json)
    for key, value in pairs(flags) do
        print(key .. " : " .. value)
    end

    local file = "---\n"
    file = file .. "title: " .. title .. "\n"
    file = file .. "date: " .. date .. "\n"
    file = file .. "---\n"

    return file
end