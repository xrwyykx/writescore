## 目录结构

~~~
|-----app //模块（为controller提供服务）
       ├─module_name            //模块目录
       │  │  ├─controller       //控制器目录（业务逻辑）
|-----bin                       //linux操作系统的脚本文件目录（如果部署到linux可以直接用）
|-----config                    //环境配置文件目录
|-----data                      //数据源目录\
|-----global                    //全局公共调用（不能引入本项目其它包）
       |------ global.go        //全局变量
       |------ config.go        //配置文件（调用config中的文件进行应用）
       |------ const.go         //全局常量
       |------ init.go          //启动入口
|------middlewares              //中间件目录
|------models                   //数据库模型
|------router                   //路由
|------utils                    //项目常用函数方法（不能引入本项目其它包）
       |------- func.go
|------ logs                    //日志                 
|------ upload                  //上传文件
|------ main.go                 //项目入口文件