package main

import (
	"context"
	"google.golang.org/grpc"
	pb "grabvn-golang-bootcamp/week03/assignment/passengerfeedback"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)



func createFeedback(ctx context.Context, c pb.PassengerFeedbackManagementClient, f *pb.PassengerFeedback) error {
	_, err := c.CreateFeedback(ctx, &pb.FeedbackReq{Feedback: f}) // inhibit result - check error only
	return err
}

func getFeedbackByPassengerID(ctx context.Context, c pb.PassengerFeedbackManagementClient, pid int32) (*pb.FeedbackList, error) {
	val, err := c.GetFeedbackByPassengerID(ctx, &pb.PassengerIDReq{PassengerID: pid})
	//log.Printf("GetFeedbackByPassengerID returns %v values", len(val.Feedbacks))
	return val, err
}

func getFeedbackByBookingCode(ctx context.Context, c pb.PassengerFeedbackManagementClient, code string) (*pb.FeedbackList, error) {
	val, err := c.GetFeedbackByBookingCode(ctx, &pb.BookingCodeReq{BookingCode: code})
	//log.Printf("GetFeedbackByBookingCode returns %v", len(val.Feedbacks))
	return val, err
}

func deleteFeedbackByPassengerID(ctx context.Context, c pb.PassengerFeedbackManagementClient, pid int32) (*pb.DeleteFeedbackByPassengerIDRes, error) {
	val, err := c.DeleteFeedbackByPassengerID(ctx, &pb.PassengerIDReq{PassengerID: pid})
	//log.Printf("DeleteFeedbackByPassengerID returns %v", val.Deleted)
	return val, err
}


func test(ctx context.Context, c pb.PassengerFeedbackManagementClient) {
	f1 := pb.PassengerFeedback{
		BookingCode: "#1",
		PassengerID: 1,
		Feedback:    "1",
	}
	err := createFeedback(ctx, c, &f1)
	if err != nil {
		panic(err)
	}
	log.Printf("Added f1: %v", &f1)

	val, err := getFeedbackByPassengerID(ctx, c, f1.PassengerID)
	if err != nil {
		panic(err)
	}
	if len(val.Feedbacks) != 1 {
		panic("getFeedbackByPassengerID returns more than one result!")
	}
	log.Println("Retrieved f1 by passengerID successfully")

	val, err = getFeedbackByBookingCode(ctx, c, f1.BookingCode)
	if err != nil {
		panic(err)
	}
	if len(val.Feedbacks) != 1 {
		panic("getFeedbackByBookingCode returns more than one result!")
	}
	log.Println("Retrieved f1 by bookingCode successfully")

	err = createFeedback(ctx, c, &f1)
	if err == nil {
		panic("Shouldn't be able to add feedback for the same booking")
	}
	log.Println("Adding f1 again failed as expected")

	err = createFeedback(ctx, c, &pb.PassengerFeedback{
		BookingCode: "#2",
		PassengerID: 1,
		Feedback:    "1",
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Added f2: %v", &f1)

	res, err := deleteFeedbackByPassengerID(ctx, c, f1.PassengerID)
	if err != nil {
		panic(err)
	}
	if res.Deleted != 2 {
		panic("deleteFeedbackByPassengerID should only return 2!")
	}
	log.Printf("Deleted all feedback related of passenger: %v", f1.PassengerID)

}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewPassengerFeedbackManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	test(ctx, c)
}
