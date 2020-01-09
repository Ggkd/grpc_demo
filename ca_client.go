package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Ggkd/v1/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func main() {
	//***********使用证书**************//


	//   生成客户端证书
	//生成的文件放在config目录下

	//	$ openssl req -new \
	//	-subj "/C=GB/L=China/O=client/CN=client.io" \
	//	-key client.key \
	//	-out client.csr

	//	$ openssl x509 -req -sha256 \
	//	-CA ca.crt -CAkey ca.key -CAcreateserial -days 3650 \
	//	-in client.csr \
	//	-out client.crt

	cert, err := tls.LoadX509KeyPair("./config/client.crt", "./config/client.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./config/ca.crt")
	if err != nil {
		fmt.Println(err)
		return
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		fmt.Println("AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "TestService",
		RootCAs:      certPool,
	})

	//建立链接
	conn, err := grpc.Dial("localhost:8886", grpc.WithTransportCredentials(c))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	//创建客户端
	client := proto.NewTestServiceClient(conn)
	//发送请求
	resp, err := client.Test(context.Background(), &proto.TestRequest{
		Request:              "client request",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	//打印响应
	fmt.Println("response------>>",resp.GetResponse())
}