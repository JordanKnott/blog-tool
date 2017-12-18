package blog

import ejson "encoding/json"
import (
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-json"
	"strings"
	"time"
)

func InitLuaRuntime(filename string) *lua.LState {
	L := lua.NewState()
	defer L.Close()

	// Register modules
	L.PreloadModule("json", json.Loader)

	// Load file and execute it
	if err := L.DoFile(filename); err != nil {
		panic(err)
	}
	return L
}

func TitleToFileName(title string) (filename string) {

	currentTime := time.Now().Local().Format("2006-01-02")
	filename = strings.Replace(title, " ", "-", -1)
	filename = strings.ToLower(filename)
	return currentTime + "-" + filename + ".md"
}

func FormatFile(L *lua.LState, title string, flags map[string]string) string {

	flagsJson, err := ejson.Marshal(flags)
	if err != nil {
		panic(err)
	}
	utcTime := time.Now().Local().Format("2006-01-02 15:04:00 -0700")
	if err := L.CallByParam(lua.P{
		Fn: L.GetGlobal("format"),
		NRet: 1,
		Protect: true,
	}, lua.LString(title), lua.LString(utcTime), lua.LString(flagsJson)); err != nil {
		panic(err)
	}

	ret := L.Get(-1)
	L.Pop(1)
	if str, ok := ret.(lua.LString); ok {
		return string(str)
	} else {
		return ""
	}
}