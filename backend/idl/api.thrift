include "./user/user.thrift"
include "./interviews/interviews.thrift"
include "./rht/mianshi.thrift"

namespace go interview


service UserService extends user.UserService {}
service InterviewsService extends interviews.InterviewsService {}
service MianshiService extends mianshi.MianshiService {}