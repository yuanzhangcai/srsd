# srsd
基于etcd的服务注册与服务发现。

服务注册example:
```
import(
    "github.com/yuanzhangcai/srsd/registry"
	"github.com/yuanzhangcai/srsd/service"
)

    //
    info := service.NewService()
    info.Name = "www.zacyuan.com"   // 服务名称
    info.Host = ":4444"             // 服务地址
    info.Metrics = ""               // prometheus指标曝露地址
    info.PProf = ""                 // pprof地址

    register = registry.NewRegistry(info,
        registry.Addresses([]string{"127.0.0.1:2379"}),
        registry.TTL(time.Duration(30*time.Second))
    err := register.Start()
    if err != nil {
        fmt.Println("服务注册失败:", err)
        return
    }

```

服务注册example:
```
import(
    "github.com/yuanzhangcai/srsd/discovery"
)

    dis = discovery.NewDiscovery(discovery.Addresses([]string{"127.0.0.1:2379"}))
    err := dis.Start("www.zacyuan.com") // 参数为空时，会自搜寻所有已注服服务。
    if err != nil {
        return err
    }

    info = dis.Select("www.zacyuan.com") // 参数为空时，从所有已注服服务信息返回其中一个服务信息。第二个参数为选择器滤器，默认为轮询过滤器。
```