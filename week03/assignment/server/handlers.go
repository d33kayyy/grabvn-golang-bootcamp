package main

import (
	"errors"
	pb "grabvn-golang-bootcamp/week03/assignment/passengerfeedback"
)

type HandlerInterface interface {
	AddFeedBack(feedback *pb.PassengerFeedback) error
	GetFeedBackByPassengerID(pid int32) ([]*pb.PassengerFeedback, error)
	GetFeedBackByBookingCode(code string) ([]*pb.PassengerFeedback, error)
	DeleteFeedBackByPassengerID(pid int32)
}

type Handler struct {
	BookingFeedbackTable   map[string][]*pb.PassengerFeedback
	PassengerFeedbackTable map[int32][]*pb.PassengerFeedback
}

func CreateHandler() Handler {
	return Handler{
		BookingFeedbackTable:   make(map[string][]*pb.PassengerFeedback),
		PassengerFeedbackTable: make(map[int32][]*pb.PassengerFeedback),
	}
}

func (h *Handler) AddFeedBack(f *pb.PassengerFeedback) (err error) {
	if _, ok := h.BookingFeedbackTable[f.BookingCode]; ok {
		err = errors.New("Feedback for this booking is already existed.")
	} else {
		h.BookingFeedbackTable[f.BookingCode] = append(h.BookingFeedbackTable[f.BookingCode], f)
		h.PassengerFeedbackTable[f.PassengerID] = append(h.PassengerFeedbackTable[f.PassengerID], f)
	}
	return
}

func (h *Handler) GetFeedBackByPassengerID(pid int32) ([]*pb.PassengerFeedback, error) {
	if val, ok := h.PassengerFeedbackTable[pid]; ok {
		return val, nil
	} else {
		return nil, errors.New("Passenger does not existed.")
	}
}

func (h *Handler) GetFeedBackByBookingCode(code string) ([]*pb.PassengerFeedback, error) {
	if val, ok := h.BookingFeedbackTable[code]; ok {
		return val, nil
	} else {
		return nil, errors.New("Booking does not existed.")
	}
}

func (h *Handler) DeleteFeedBackByPassengerID(pid int32) (deleted int32) {
	if val, ok := h.PassengerFeedbackTable[pid]; ok {
		deleted = int32(len(val))
		for _, f := range val {
			delete(h.BookingFeedbackTable, f.BookingCode)
		}

		delete(h.PassengerFeedbackTable, pid)

	}
	return
}
