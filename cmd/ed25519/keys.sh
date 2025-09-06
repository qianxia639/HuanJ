#!/bin/bash

# 创建临时文件
TMP_FILE=$(mktemp /tmp/temp_keys.XXXXXX.go) || exit 1

# 确保临时文件被删除
trap 'rm -f "$TMP_FILE"' EXIT

cat > $TMP_FILE << EOF
package main

import "fmt"

func main() {
    fmt.Println("Hello World!")
}
EOF

go  run $TMP_FILE

# rm -f /tmp/temp_keys.go