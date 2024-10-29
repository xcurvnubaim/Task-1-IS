package shareRequest

type RequestStatusEnum struct {
	Accepted string
	Rejected string
	Pending string
}

var RequestStatus = &RequestStatusEnum{
	Accepted: "accepted",
	Rejected: "rejected",
	Pending: "pending",
}