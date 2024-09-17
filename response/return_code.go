package response

// RetCode type for return codes
type RetCode string

// Define common ret codes for different scenarios
const (
    // Success Codes
    SuccessOK              RetCode = "200"
    SuccessCreated         RetCode = "201"
    
    // Client Error Codes
    BadRequest             RetCode = "400"
    Unauthorized           RetCode = "401"
    Forbidden              RetCode = "403"
    NotFound               RetCode = "404"
    UnprocessableEntity    RetCode = "422"
    
    // Server Error Codes
    InternalServerError    RetCode = "500"
    ServiceUnavailable     RetCode = "503"
)
