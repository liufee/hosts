Hosts管理
===============================

基于gui库[fyne](https://github.com/fyne-io/fyne)的桌面管理hosts文件程序，支持mac os、windows、linux


编译环境
-------
 * Go版本 >= 1.12
 * CGO
详细参考[fyne](https://github.com/fyne-io/fyne)


安装
-------
 1. 下载可执行文件  
   [点击进入下载](https://github.com/liufee/hosts/releases) 
   
   
 2. 编译安装
  ```bash 
        $ git clone https://github.com/liufee/hosts.git
        $ cd hosts
        $ sh build.sh
  ```
    

运行
-------
windows需要以管理员身份运行(操作c:\windows\system32\drivers\etc\hosts文件)  
macos、linux需要保证运行的用户能读写/etc/hosts文件