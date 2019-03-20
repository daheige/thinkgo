# ubuntu mysql5.7最大连接数设置
    修改mysql中的mysqld.cnf文件
    cd /etc/mysql/mysql.conf.d
    执行下面的操作来修改配置文件
    sudo vim mysqld.cnf
    在 [mysqld] 中新加
    max_connections  =1000；
    按下esc按键输入：wq保存退出

    重启mysql服务器
    sudo service mysql restart；

    登录进去查看mysql配置
    mysql -uroot -p
    //输入密码登录
    //查看刚刚配置信息
    mysql>show variables like '%max_connections%';
    发现没有改为默认最大是214
    因为ubuntu系统本身有限制文件打开和连接数量，所以需要修改系统配置来达到我们的要求

    修改系统配置
    cd  /etc/systemd/system/multi-user.target.wants
    sudo vim mysql.service
    //在 [Service] 最后加入：
    LimitNOFILE=65535
    LimitNPROC=65535
    按下esc按键输入：wq保存退出

    刷新系统配置
    systemctl daemon-reload
    systemctl restart mysql.service

    检验配置是否成功
    mysql -uroot -p
    //输入密码登录
    //查看刚刚配置信息
    mysql>show variables like '%max_connections%';

    这是可以看到已修改为1000了
