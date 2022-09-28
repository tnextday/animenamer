# animenamer

animenamer 是一个动画（anime）重命名工具，它可以使用`集数`（AbsoluteNumber）或者
`季集`（season&episode）来从 [thetvdb](https://thetvdb.com) 或 [The Movie Database (TMDB)](https://www.themoviedb.org/)上获取信息，来重命名文件，使其文件名对搜刮器更友好。

## 主要特性

* 可以使用`集数`或`季集`来索引信息
* 可以使用 [TVDB](https://thetvdb.com) 或者 [TMDB](https://www.themoviedb.org/) 来索引信息
* 使用正则表达式匹配文件
* 可定制文件名
* 重命名字幕文件
* 文件整理，可以将搜索到的文件移动到指定文件夹
* 显示缺失的剧集
* 单文件无依赖，下载即可运行
* 自定义剧集信息翻译
* ~~生成 kodi nfo 文件~~
* ~~下载剧集图片~~

_v0.4 版本以后去掉了生成 kodi nfo 和下载图片的功能_

## 安装

#### 二进制

https://github.com/tnextday/animenamer/releases

## 使用方法

#### 基本使用

```
animenamer -n 海贼王 -p '.*第(?P<absolute>\d+)集\.mkv' -d ./
animenamer -n 海贼王 -p '.*第(?P<absolute>\d+)集\.mkv' ./
```
参数说明:

- `-n`: 指定剧集名称
- `-p`: 匹配文件名的正则表达式，同时指定 `absolute` （集数）在文件名中的位置
- `-d`: 只打印重命名文件，不真正执行
- `./`: 文件所在目录

#### 使用配置文件

animenamer 支持使用配置文件来设置相关参数。animenamer 默认会使用当前目录下
`animenamer.yaml`作为配置文件，也可以使用 `-c` 参数指定配置文件。

在当前目录写入默认配置文件

```
animenamer writeConfig
```

参考 [animenamer.yml](examlpes/animenamer.yml)

另外可以使用配置文件`animenamer.custom.yml`覆盖tvdb的资料信息，例如名字、简介、季名等，参考 [animenamer.custom.yml](examlpes/animenamer.custom.yml)

#### 详细帮助
```
animenamer -h
```

#### 代理

使用代理服务器
```
export HTTP_PROXY=<proxy-url>
export HTTPS_PROXY=<proxy-url>
animenamer [flags] <anime-files>
```

## 正则表达式匹配与新文件名格式

*正则表达式使用的是 re2 语法，详细语法请访问 [re2 Syntax](https://github.com/google/re2/wiki/Syntax)*。

**[regular expression tester](https://regoio.herokuapp.com/)**

#### 正则表达式参数

animenamer 使用正则表达式中的`命名分组`来提取文件名中的剧集信息，命名分组的格式为`(?P<name>re)`

表达式中**必须包含 “`absolute`” 或者 “`season`和`episode`” 命名**。
此外，如果命令行参数或者配置中没指定`series`或者`seriesId`参数，
那么也需要在正则表达式中指定。

例如：

```
正则：.*第(?P<absolute>\d+)集\.mkv
文件名：海贼王第10集.mkv
```

匹配到的集数为 **10**

#### 新文件名格式

新文件名使用 `{variable}` 格式来替换相关变量，另外支持 `{variable.n}` 这种格式对数字添 0 补位

除了使用系统保留变量 （`series`, `seriesId`, `season`, `episode`, `absolute`, `date`, `title`, `ext`），
也可以使用命名分组匹配到的变量

例如

```
正则: .*\.(?P<absolute>\d+)\.(?P<codec>.+)\.mp4
文件名: op.001.x264.mp4
新文件名格式：{series}.S{season.2}E{episode.2}.[{absolute.3}].{codec}.{ext}
```
那么新文件名为 `海贼王.S01E01.[001].x264.mp4`


## 编译

```
make
```