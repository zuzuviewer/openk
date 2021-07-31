# Openk

json转yaml，yaml转json，输出格式化和字符串格式的json数据工具，从指定文件中得到输入。
支持将结果输出到标准输出或者指定文件

# Install

go get https://github.com/zuzuviewer/openk

cd openk/cmd

go install openk.go

# Usage

## json转yaml

openk convert -f test.json

openk convert -f test.json -w test.yaml

## yaml转json

openk convert -f test.yaml

openk convert -f test.yml -w test.json

## 输出格式化的json

openk convert -i -f test.json

openk convert -i -f test.json -w t.json

## 输出字符串格式的json

openk convert -s -f test.json

openk convert -s -f test.json -w t.txt

# License

Openk is released under the Apache 3.0 license.See [LICENSE](./LICENSE)
