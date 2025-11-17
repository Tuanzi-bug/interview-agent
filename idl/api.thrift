include "./user/user.thrift"
include "./interviews/interviews.thrift"
include "./titleBank/titleBank.thrift"

namespace go interview


service UserService extends user.UserService {}
service InterviewsService extends interviews.InterviewsService {}
service TitleBankService extends titleBank.TitleBankService {}
