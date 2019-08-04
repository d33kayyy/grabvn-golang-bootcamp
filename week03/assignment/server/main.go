package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grabvn-golang-bootcamp/week03/assignment/passengerfeedback"
)

const (
	port = ":50051"
)

// server is used to implement
type server struct{}

var handler = CreateHandler()

func (s *server) CreateFeedback(ctx context.Context, in *pb.FeedbackReq) (*pb.FeedbackRes, error) {
	log.Printf("CreateFeedback Received: %v", in.Feedback)
	err := handler.AddFeedBack(in.Feedback)
	res := &pb.FeedbackRes{Feedback: in.Feedback}
	return res, err
}

func (s *server) GetFeedbackByPassengerID(c context.Context, in *pb.PassengerIDReq) (*pb.FeedbackList, error) {
	log.Printf("GetFeedbackByPassengerID Received: %v", in.PassengerID)
	l, err := handler.GetFeedBackByPassengerID(in.PassengerID)
	return &pb.FeedbackList{Feedbacks: l}, err
}

func (s *server) GetFeedbackByBookingCode(c context.Context, in *pb.BookingCodeReq) (*pb.FeedbackList, error) {
	log.Printf("GetFeedbackByBookingCode Received: %v", in.BookingCode)
	l, err := handler.GetFeedBackByBookingCode(in.BookingCode)
	return &pb.FeedbackList{Feedbacks: l}, err
}

func (s *server) DeleteFeedbackByPassengerID(c context.Context, in *pb.PassengerIDReq) (*pb.DeleteFeedbackByPassengerIDRes, error) {
	log.Printf("DeleteFeedbackByPassengerID Received: %v", in.PassengerID)
	deleted := handler.DeleteFeedBackByPassengerID(in.PassengerID)
	return &pb.DeleteFeedbackByPassengerIDRes{Deleted: deleted}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPassengerFeedbackManagementServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
