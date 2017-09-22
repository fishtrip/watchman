# Golang 使用 Capistrano 部署配置

------------------------------

本目录用于 Golang 项目首次配置自动部署使用

## 初始化项目部署配置

1. 拷贝本目录(source/golang/)的 **所有文件** 到Golang 项目的 **根目录** 下
2. 修改部分文件的配置 (以下文件路径以项目中的文件路径为准)
  1. 修改 `config/deploy.rb` 文件

  ```ruby
    
    # config/deploy.rb
    # 修改项目名称，例如
    set :application, "jarvis"
    # 修改项目的git地址，例如
    set :repo_url, 'git@gitlab.fishtrip.cn:ruby/jarvis.git'
    # 设置项目在生产服务器的部署目录，通常放在 /opt/work/ 下，以项目名为部署目录，例如
    set :deploy_to, '/opt/work/jarvis'
    
    ...
    ...
    
    # 项目启动服务的参数如果需要修改，在文件的最下方：
    execute "sh", "restart.sh"
  ```
  2. 配置 `build.rb`

  ```ruby
  # build.rb
  
  application = "jarvis"
  ```
  
  3. 配置测试和线上的部署服务器

  ```ruby
    # config/deploy/alpha.rb  测试部署服务器配置 
    
    # 修改测试服务器，通常是 a3
    server 'a3', user: 'deploy', roles: %w{app}
  ```
  
  ```ruby
    # config/deploy/production.rb  生产环境部署服务器配置 
    
    # 修改生产服务器，如果需要部署多台，写多行配置即可
    server 'r5', user: 'deploy', roles: %w{app}
    server 'r6', user: 'deploy', roles: %w{app}
  ```
  
  
## 测试／生产环境部署

1. 登陆gitlab，进入项目目录后，点击右上角 `⚙️` -> `Deploy Keys`， 分别找到 `jump` / 测试服务器的ssh key / 生产环境要部署的服务器 ssh key， 点击 `enable` （**注：未激活对应部署环境的deploy key，项目无法部署**）
1. 访问 **跳板服务器** ，将项目代码 clone 到 `/opt/work/项目目录` 下
1. 进入项目目录，执行 `bundle install`
1. 部署测试服务器执行 `bundle exec cap alpha deploy`，部署生产环境执行  `bundle exec cap production deploy`
1. 若命令成功执行完成没有报错，说明部署成功


## 补充
1. 项目部署完成后，确认是否项目启动成功，去运行服务器部署目录的/run 目录查看是否有 `application.pid` 文件(例如avatar项目的pid文件，全路径是 `/opt/work/avatar/current/run`)，存在即启动成功
1. 项目的启动端口需要在maven编译的配置文件里，具体的使用端口，请在部署前找运维确认分配的端口

