include "./user/user.thrift"
include "./interviews/interviews.thrift"
include "./demo/demo.thrift"

namespace go interview


service UserService extends user.UserService {}
service InterviewsService extends interviews.InterviewsService {}
service DemoService extends demo.DemoService {}