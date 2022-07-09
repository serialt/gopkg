## go函数库

### 使用方法

```
    go get -u  github.com/serialt/gopkg
```



### 函数清单

`shell`

| 函数名                      | 说明                            | 备注 |
| --------------------------- | ------------------------------- | ---- |
| FindCommandPath,Where       | 查找命令的路径                  |      |
| RunCmd                      | 获取标准正确输出                |      |
| RunCommandWithTimeout       | 带超时控制的执行shell命令       |      |
| FindUser                    | 查找操作系统上的用户            |      |
| GetLoginUser,GetCurrentUser | 获取当前登录的用户              |      |
| UserHomeDir                 | 获取当前用户家目录              |      |
| UserCacheDir                | 当前用户缓存目录，$HOME/.cache  |      |
| UserConfigDir               | 当前用户配置目录，$HOME/.config |      |
| Hostname                    | 获取主机名                      |      |
| IsMSys                      | 判断是否是msys(MINGW64)         |      |
| CurrentShell                | 获取操作系统的shell             |      |
| HasShellEnv                 | 判断操作系统是否有此shell       |      |
| IsShellSpecialVar           | 判断是否是shell的变量字符       |      |
| Workdir                     | 获取workspace                   |      |

`string`

| 函数名                      | **说明**                                        | **备注** |
| --------------------------- | ----------------------------------------------- | -------- |
| StringIsEmpty               | 判断字符串是否为空，是则返回true，否则返回false |          |
| StringIsNotEmpty            | 判断字符串是否不为空                            |          |
| StringConvert               | 去掉下划线并将下划线后的首字母大写              |          |
| StringRandSeq               | 创建指定长度的随机字符串                        |          |
| StringRandSeq16             | 创建长度为16的随机字符串                        |          |
| StringAllLetter             | 判断字符串是否只由字母组成                      |          |
| StringTrim                  | 去除字符串中的空格和换行符                      |          |
| StringTrimN                 | 去除字符串中的换行符                            |          |
| ToString                    | 将对象格式化成字符串                            |          |
| StringSingleValue           | 将字符串内所有连续value替换为单个value          |          |
| StringSingleSpace           | 将字符串内所有连续空格替换为单个空格            |          |
| StringPrefixSupplementZero  | 当字符串长度不满足时，将字符串前几位补充0       |          |
| SubString                   | 截取字符串                                      |          |
| StringBuild，StringBuildSep | 拼接字符串                                      |          |
| FilterPrefix                | 根据前缀过滤slice                               |          |
| FindLongestStr              | 查询最长字符串                                  |          |
| ArrayToString               | 数字切片变字符串                                |          |
| StructToMap                 | 结构体转map                                     |          |

`time`

| 函数名           | **说明**         | **备注** |
| ---------------- | ---------------- | -------- |
| Timestamp2String | 时间戳转字符串   |          |
| String2Timestamp | 字符串转时间戳   |          |
| GetDate          | 返回系统当前时间 |          |
| GetRunTime       | 获取当前系统环境 |          |

`file`

| 函数名             | **说明**                                                     | **备注** |
| ------------------ | ------------------------------------------------------------ | -------- |
| GetRootPath        | 获取项目路径                                                 |          |
| DirExist，IsDir    | 判断目录否存在                                               |          |
| Mode               | unix类系统获取文件的权限                                     |          |
| FileExt，Suffix    | 获取文件的后缀, main.go 获取的后缀是.go                      |          |
| Prefix             | 获取文件名前缀, /tmp/main.go 获取的文件前缀是main            |          |
| CreateEmptyFile    | 创建空文件                                                   |          |
| FileExists，IsFile | 文件是否存在                                                 |          |
| IsAbsPath          | 是否是绝对路径                                               |          |
| PathDir            | 获取路径的目录                                               |          |
| Name               | 获取路径的文件名                                             |          |
| DeleteFile         | 删除文件或目录                                               |          |
| MkDir              | 创建文件夹,支持x/a/a  多层级                                 |          |
| FilePathExists     | 判断路径是否存在                                             |          |
| FileReadFirstLine  | 从文件中读取第一行并返回字符串数组                           |          |
| FileReadPointLine  | 从文件中读取指定行并返回字符串数组                           |          |
| FileReadLines      | 从文件中逐行读取并返回字符串数组                             |          |
| FileParentPath     | 文件父路径                                                   |          |
| ReadFile           | 读文件                                                       |          |
| WriteFile          | 写文件                                                       |          |
| WriteStringToFile  | 带权限位写文件                                               |          |
| FileAppend         | 文件追加内容                                                 |          |
| RecreateDir        | 创建目录                                                     |          |
| GetFilepaths       | 获取目录里的所有文件                                         |          |
| GetFiles           | 获取文件，返回文件路径和内容                                 |          |
| FileLoopDirs       | 遍历目录下的所有子目录，即返回pathname下面的所有目录，目录为绝对路径 |          |
| FileLoopOneDirs    | 遍历目录下的所有子目录，即返回pathname下面的所有目录，目录为相对路径 |          |
| FileLoopFileNames  | 遍历文件夹及子文件夹下的所有文件名，即返回pathname目录下所有的文件，文件名为相对路径 |          |
| FileMove           | 移动文件                                                     |          |
| TrimSpace          | 去除空格                                                     |          |
| FileCompressZip    | zip压缩文件                                                  |          |
| FileDeCompressZip  | zip解压文件                                                  |          |
| FileCompressTar    | tar压缩文件                                                  |          |
| FileCopy           | 文件复制                                                     |          |
| IsImageFile        | 判断是否是图片                                               |          |
| IsZipFile          | 判断是否zip压缩文件                                          |          |
| OSTempFile         | 创建临时文件                                                 |          |
| TempFile           | 指定目录里创建临时文件                                       |          |
| UserHomePath       | 用户家目录                                                   |          |
| Mkdir              | 带权限创建目录                                               |          |
| OpenFile           | 打开一个文件                                                 |          |
| QuickOpenFile      | 快速打开一个文件                                             |          |
| CreateFile         | 创建文件                                                     |          |
| DeleteIfExist      | 删除文件                                                     |          |
| Unzip              | 解压zip                                                      |          |

`byte`

| 函数名        | 说明             | 备注 |
| ------------- | ---------------- | ---- |
| GetBytes      | 获取接口字节数组 |      |
| IntToBytes    | int转换成字节    |      |
| BytesToInt    | 字节转换成int    |      |
| Uint16ToBytes | uint16转换成字节 |      |
| BytesToUint16 | 字节转换成uint16 |      |
| Uint32ToBytes | uint32转换成字节 |      |
| BytesToUint32 | 字节转换成uint32 |      |
| Uint64ToBytes | uint64转换成字节 |      |
| BytesToUint64 | 字节转换成uint64 |      |

`env`

| 函数名         | 说明                            | 备注 |
| -------------- | ------------------------------- | ---- |
| EnvGet         | 环境变量名称                    |      |
| EnvGetD        | 环境变量为空时的默认值          |      |
| EnvGetInt      | 获取环境变量 envName 的值       |      |
| EnvGetIntD     | 环境变量为空时的默认值          |      |
| EnvGetInt64    | 获取环境变量的值                |      |
| EnvGetInt64D   | 环境变量为空时的默认值          |      |
| EnvGetUint64   | 获取环境变量的值                |      |
| EnvGetUint64D  | 获取环境变量的值                |      |
| EnvGetFloat64  | 获取环境变量的值                |      |
| EnvGetFloat64D | 获取环境变量的值                |      |
| EnvGetBool     | 获取环境变量的值                |      |
| Environ        | 获取所有的环境变量，返回map数据 |      |
| IsAIX          | 判断操作系统类型                |      |
| IsAndroid      | 判断操作系统类型                |      |
| IsMac          | 判断操作系统类型                |      |
| IsDarwin       | 判断操作系统类型                |      |
| IsFreeBSD      | 判断操作系统类型                |      |
| IsIOS          | 判断操作系统类型                |      |
| IsLinux        | 判断操作系统类型                |      |
| IsNetBSD       | 判断操作系统类型                |      |
| IsOpenBSD      | 判断操作系统类型                |      |
| IsPlan9        | 判断操作系统类型                |      |
| IsWin          | 判断操作系统类型                |      |
| IsWindows      | 判断操作系统类型                |      |
| Is386          | 判断操作系统类型                |      |
| IsAMD64        | 判断操作系统类型                |      |
| IsRISCV64      | 判断操作系统类型                |      |

`excelize`

| 函数名          | 说明                    | 备注 |
| --------------- | ----------------------- | ---- |
| PasreExcel2List | 读取excel文件到二位切片 |      |
| PasreList2Excel | 解析二位切片到excel文件 |      |

`hash`

| 函数名          | 说明 | 备注 |
| --------------- | ---- | ---- |
| HashMD5Bytes    |      |      |
| HashMD5         |      |      |
| HashMD516Bytes  |      |      |
| HashMD516       |      |      |
| HashSha1Bytes   |      |      |
| HashSha1        |      |      |
| HashSha224Bytes |      |      |
| HashSha224      |      |      |
| HashSha256Bytes |      |      |
| HashSha256      |      |      |
| HashSha384Bytes |      |      |
| HashSha384      |      |      |
| HashSha512Bytes |      |      |
| HashSha512      |      |      |

`ip`

| 函数名          | 说明                                                         | 备注 |
| --------------- | ------------------------------------------------------------ | ---- |
| IPGet           | 返回客户端 IP                                                |      |
| GetPubIP        | 获取公网ip, 如果两个ip不同，则访问ip.tool.lu 和false, ip都相同则返回true |      |
| SubNetMaskToLen | ipv4 子网掩码长度换算                                        |      |
| LenToSubNetMask | ipv4 网络位长度转换为子网掩码地址                            |      |
| IsPublicIPv4    | ipv4 判断是否是公网ip                                        |      |











