include "./user/user.thrift"
include "./interviews/interviews.thrift"
include "./questionbank/questionbank.thrift"

namespace go interview


service UserService extends user.UserService {}
service InterviewsService extends interviews.InterviewsService {}
service QuestionBankService extends questionbank.QuestionBankService {}
