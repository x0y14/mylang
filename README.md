# MyLang
自作言語

## compiler
自作言語を解析して，runtime用のオレオレアセンブリを生成します．

## runtime
オレオレアセンブリを読み込んで動きます．  
[runtime/runtime_test.go](runtime/runtime_test.go)に`TestRuntime_Run_FizzBuzz`関数があるので，動作が気になる方はこれをチェックしてください．