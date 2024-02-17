





# 一.函数原型

## AddCommand

```C#
public static void AddCommand(
    string lpName, 
    string filePath, 
    string loadType, 
    string description,
    string usage
)
```

## AddMenuItemA

```C#
 public void AddMenuItemA(
     string menuName, 
     string iconPath
 )
```



## AddMenuItemB

```C#
public void AddMenuItemB(
    string menuName, 
    string iconPath, 
    NLua.LuaFunction clickEvent
)
```



## AttackMenuExpand 

```C#
  public void AttackMenuExpand(
      string menuName,
      string iconPath, 
      NLua.LuaFunction clickEvent  
  )
```





## AddTextBox

```C#
TextBox AddTextBox(
    Form form,
    int x,
    int y, 
    int width, 
    int height
)
```



## AddSelectFileButton

```C#
 public Button AddSelectFileButton(
     Form form,
     string text,
     int x, 
     int y, 
     int width,
     int height
 )
```

## Addlabel

```C#
public Label Addlabel(
    Form form, 
    string title, 
    int x, 
    int y, 
    int width, 
    int height
)
```

## AddTextBox

```C#
public TextBox AddTextBox(
    Form form, 
    int x, 
    int y, 
    int width, 
    int height
)
```

## AddComBox

```C#
public ComboBox AddComBox(
    Form form,
    int x, 
    int y,
    int width,
    int height
)
    
    
```



## AddButton

```C#
 public void AddButton(
     Form form, 
     string text, 
     int x, 
     int y, 
     int width,
     int height, 
     NLua.LuaFunction clickEvent
 )

```



## AddMenuItemsAsSubItems

```C#
 public void AddMenuItemsAsSubItems(
     string parentItemName, 
     LuaTable subItemsTable
 )
```



## CreateForm

```C#
 public Form CreateForm(
     string title,
     int width,
     int height
 )
```

## ExecuteAssembly

```C#
public static void ExecuteAssembly(
    string filePath, 
    string args
)
```



## Inlineassembly

```C#
 public static void Inlineassembly(
     string filePath, 
     string args
 )
```



## GetFileName

```C#
public string GetFileName(string filePath)
```



## MessageboxA

```C#
public void MessageboxA(string conText)
```



## MenuStripExpand

```C#
public void MenuStripExpand(  
    string menuName,
    string iconPath,
    NLua.LuaFunction clickEvent   //Click Event
)
```



## Nopowershell

```C#
public static void Nopowershell(
    string command, 
    string outString
)
```



## Upload

```C#
public static void Upload(
    string uploadFilePath, 
    string filePath
)
```



## PEloader

```C#
public static void PEloader(
    string filePath, 
    string args
)
```

## Sessionlog

```C#
public void Sessionlog(string conText)
```





# 二.lua  Demo

lua Plugins demo

```lua
local function SchTaskForm()
    local IsOK = false
    local Form1 = CreateForm("sch", 422, 355)

    local Label1 = Addlabel(Form1, "Interval(min):", 17, 33, 133, 25)
    local Label2 = Addlabel(Form1, "Select File:", 31, 103, 114, 25)
    local Label3 = Addlabel(Form1, "Directory:", 28, 174, 114, 25)

    local addButton = AddSelectFileButton(Form1, "Select", 162, 100, 181, 33)
    local texBox1 = AddTextBox(Form1, 162, 31, 181, 33)
    local ComBox1 = AddComBox(Form1, 162, 174, 181, 33)
    local Button1 = AddButton(Form1, "Ok", 36, 234, 106, 37, function()
        if texBox1.Text ~= "" and ComBox1.Text ~= "" and addButton.Text ~=
            "Select" then -- The values of texBox1 and ComBox1 can't be empty strings, otherwise this window
            IsOK = true
            Form1:Hide()

        else
            IsOK = false
            MessageBoxA("adas")
        end
    end)

    local Button2 = AddButton(Form1, "close", 236, 234, 107, 37,
                              function() Form1:Close() end)
    Form1:ShowDialog()
    if IsOK then
        local FileName = GetFileName(addButton.Text)
        local UploadFilePath = ComBox1.Text .. "\\" .. FileName
        Upload(UploadFilePath, addButton.Text)
        -- InlineAssembly("Plugins\\BypassUAC-ETV.exe", "Plugins\\BypassUAC-ETV.exe", texBox1.Text)
    end
end


AddCommand(
    "SharpKatz",
    "Plugins\\SharpKatz.exe", --FilePath
    "execute-assembly",
    "Steal domain login credentials",
    "SharpKatz -h"
)
AddCommand(
    "fscan",
    "Plugins\\fscan.bin",
    "RunPE",
    "asdasd",
    "DDDDD"
)
AttackMenuExpand("asdasd","",null)
MenuStripExpand("asdasd",null,null)


local function Privileg()
    InlineAssembly("Plugins\\BypassUAC.exe", "calc")
end

local function RunPEload()
    RunPE("Plugins\\fscan.bin","-h 192.168.1.1/24")
    -- 调用C#中注册的printf方法
    Sessionlog("Message")
end

local function AddMenuItem()

    AddWindowsMenuC("渗透插件", null)
    AddWindowsMenuC("信息收集", null)
    AddWindowsMenuC("抓取浏览器密码", null)
    AddWindowsMenuC("持久化控制", null)
    AddWindowsMenuC("任务计划", null)
    AddWindowsMenuD("隐藏安装", "", SchTaskForm)
    AddWindowsMenuD("权限提升",null,function ()
        Sessionlog("Executing program..")
    end)
    AddWindowsMenuD("4.0","",RunPEload)
    local menuStructure = {
        ["渗透插件"] = {"信息收集", "持久化控制","权限提升"},
        ["信息收集"] = {"抓取浏览器密码"},
        ["持久化控制"] = {"任务计划"},
        ["任务计划"] = {"隐藏安装"},
        ["权限提升"] = {"4.0"}
    }
    for parent, subs in pairs(menuStructure) do
        AddMenuItemsAsSubItems(parent, subs)
    end
end


AddMenuItem()

```

