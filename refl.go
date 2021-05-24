package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

func main() {
	address := "localhost:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	client := rpb.NewServerReflectionClient(conn)

	if err := singleStream(client); err != nil {
		log.Fatal(err)
	}
}

func singleStream(client rpb.ServerReflectionClient) error {
	fmt.Println("*** single stream 💥 ***")
	stream, err := client.ServerReflectionInfo(context.Background())
	if err != nil {
		return err
	}

	req := &rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_ListServices{},
	}
	if err := stream.Send(req); err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	write(req, resp)
	symbols := resp.GetListServicesResponse().GetService()
	for _, symbol := range symbols {
		req := &rpb.ServerReflectionRequest{
			MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: symbol.GetName(),
			},
		}
		if err := stream.Send(req); err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		write(req, resp)
	}

	fmt.Println("waiting for another response:")
	resp, err = stream.Recv()
	if err != nil {
		return err
	}
	write(req, resp)

	return nil
}

func write(req *rpb.ServerReflectionRequest, resp *rpb.ServerReflectionResponse) {
	fmt.Println("Req                 : ", req)
	fmt.Println("Resp.OriginalRequest: ", resp.GetOriginalRequest())
	//fmt.Println("Resp (full)         : ", resp)
	fmt.Println()

}

func mutipleStreams(client rpb.ServerReflectionClient) error {
	fmt.Println("*** multiple streams ✅ ***")
	stream, err := client.ServerReflectionInfo(context.Background())
	if err != nil {
		return err
	}

	req := &rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_ListServices{},
	}
	if err := stream.Send(req); err != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	write(req, resp)
	symbols := resp.GetListServicesResponse().GetService()
	for _, symbol := range symbols {
		stream, err := client.ServerReflectionInfo(context.Background())
		if err != nil {
			return err
		}

		req := &rpb.ServerReflectionRequest{
			MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{
				FileContainingSymbol: symbol.GetName(),
			},
		}
		if err := stream.Send(req); err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		write(req, resp)
	}
	return nil
}

func callList(client rpb.ServerReflectionClient, cnt int) error {
	fmt.Println("*** call list ***")
	stream, err := client.ServerReflectionInfo(context.Background())
	if err != nil {
		return err
	}
	req := &rpb.ServerReflectionRequest{MessageRequest: &rpb.ServerReflectionRequest_ListServices{}}
	for i := 0; i < cnt; i++ {
		stream.Send(req)
		stream.Recv()
	}
	return nil
}
