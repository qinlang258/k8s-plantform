# k8s-plantform
k8s-plantform k8s平台的练习项目
主要思路

1 首先在dataselector创建一个 想要的资源cells 例如 podCell corev1.Pod  
2 创建service资源，使用k8s的client-go获取资源，然后根据dataSelector 转化资源类型为自定义资源  
3 创建controller 根据请求的路由将请求转发到k8s，再由dataSelector转换  
