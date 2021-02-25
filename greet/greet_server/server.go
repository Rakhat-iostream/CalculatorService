package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"rakhat/greet/greetpb"
	"time"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}
type CalculatorService struct {
	greetpb.UnimplementedCalculatorServiceServer
}

func (s CalculatorService) GetSum(ctx context.Context, req *greetpb.NumbersRequest) (*greetpb.NumbersResponse, error) {
	fmt.Printf("Sum function was invoked with %v \n", req)
	first, second := req.GetFirstNumber(), req.GetSecondNumber()
	return &greetpb.NumbersResponse{Result: first + second}, nil
}
func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *CalculatorService) PrimeNumberDecomposition(req *greetpb.IntegerRequest, stream greetpb.CalculatorService_PrimeNumberDecompositionServer) error {
	num := req.GetNumber()
	primes := primeNumberDecomposedResult(num)
	for i := 0; i < len(primes); i++ {
		res := &greetpb.IntegerResponse{Result: primes[i]}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error with responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}

func primeNumberDecomposedResult(n int32) []int32 {
	arr := []int32{}
	for {
		if n%2 != 0 {
			break
		}
		arr = append(arr, 2)
		n /= 2
	}
	var i int32 = 0
	for i = 3; i <= n*n; i += 2 {
		for {
			if n%i != 0 {
				break
			}
			arr = append(arr, i)
			n /= i
		}
	}
	if n > 2 {
		arr = append(arr, n)
	}

	return arr
}

func (s *CalculatorService) ComputeAverage(stream greetpb.CalculatorService_ComputeAverageServer) error {
	var sum int32 = 0
	var count int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &greetpb.AverageResponse{Result: float64(sum / count)}
			return stream.SendAndClose(response)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetNumber()
		count++
	}
}
func main() {
	fmt.Printf("Listening on port 50051")
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterCalculatorServiceServer(s, &CalculatorService{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
