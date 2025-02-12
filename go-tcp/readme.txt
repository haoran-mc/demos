

组装一个 tcp 报文并发送。

    sudo go run tcp.go



使用 tcpdump 监听回环网卡，查看发送的报文。

    sudo tcpdump -i lo0 port 9022
