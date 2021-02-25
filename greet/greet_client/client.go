package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"rakhat/greet/greetpb"
	"time"
)

func main() {
	fmt.Println("Hello, i'm client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldn't connect: %v", err)
	}
	defer conn.Close()
	c := greetpb.NewCalculatorServiceClient(conn)
	//PrintPrimeResponses(120, c)

	PrintAverage([]int32{1, 2, 3, 4, 5}, c)
}

func PrintPrimeResponses(n int32, c greetpb.CalculatorServiceClient) {
	req := &greetpb.IntegerRequest{Number: n}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("error with server stream RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("error with response from server stream RPC %v", err)
		}
		log.Printf(fmt.Sprint(res.GetResult(), " "))
	}
}

func PrintAverage(numbers []int32, c greetpb.CalculatorServiceClient) {
	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}
	for _, n := range numbers {
		stream.Send(&greetpb.IntegerRequest{Number: n})
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from server ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res.GetResult())
}
func Sum(c greetpb.CalculatorServiceClient, first int32, second int32) {
	ctx := context.Background()
	request := &greetpb.NumbersRequest{
		FirstNumber:  first,
		SecondNumber: second,
	}
	response, err := c.GetSum(ctx, request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC $v", err)
	}
	log.Printf("response from CalculatorService: %v", response.Result)
}
