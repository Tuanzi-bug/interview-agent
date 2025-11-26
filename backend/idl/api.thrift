include "./user/user.thrift"
include "./interviews/interviews.thrift"

namespace go interview


service UserService extends user.UserService {}
service InterviewsService extends interviews.InterviewsService {}