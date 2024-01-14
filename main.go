package main

import (
	"fmt"
	"miniKV/db"
)

func main() {
	d := openHandler()

	opt := ""
	for {
		fmt.Printf("请输入指令, 输入-h查看帮助\n")
		fmt.Scanf("%s\n", &opt)
		fmt.Printf("opt: %s\n", opt)
		switch opt {
		case "get":
			getHandler(d)
		case "put":
			putHandler(d)
		case "del":
			delHandler(d)
		case "merge":
			mergeHandler(d)
		case "-h":
			fmt.Printf("get put del merge\n")
		case "exit":
			d.Close()
			fmt.Printf("退出")
			return
		}
		opt = ""
	}
}

func mergeHandler(d *db.DB) {
	err := d.Merge()
	if err != nil {
		fmt.Printf("merge err: %+v\n", err)
	}
}

func openHandler() *db.DB {
	name := ""
	fmt.Printf("请输入数据库名称\n")
	fmt.Scanf("%s\n", &name)
	d, err := db.OpenDB(name)
	if err != nil {
		fmt.Printf("连接数据库出错: %+v\n", err)
		return nil
	}
	fmt.Printf("%+v\n", d)
	return d
}

func putHandler(d *db.DB) {
	k, v := "", ""
	fmt.Scanf("%s %s\n", &k, &v)
	if k == "" || v == "" {
		fmt.Printf("key或value不能为空")
		return
	}
	err := d.Put(k, v)
	if err != nil {
		fmt.Printf("插入出错: %+v\n", err)
		return
	}

	fmt.Printf("插入成功\n")

}

func delHandler(d *db.DB) {
	k := ""
	fmt.Scanf("%s\n", &k)
	if k == "" {
		fmt.Printf("key不能为空")
		return
	}
	err := d.Del(k)
	if err != nil {
		fmt.Printf("删除出错: %+v\n", err)
		return
	}

	fmt.Printf("删除成功\n")

}

func getHandler(d *db.DB) {
	k := ""
	fmt.Scanf("%s\n", &k)
	if k == "" {
		fmt.Printf("key不能为空")
		return
	}
	val, ok, err := d.Get(k)
	if err != nil {
		fmt.Printf("查询出错: %+v\n", err)
		return
	}
	if !ok {
		fmt.Printf("key: %s不存在\n", k)
		return
	}

	fmt.Printf("value: %s\n", val)

}
