package v3

type WusthelperResp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type StudentInfoResp struct {
	StuNum      string `json:"stuNum"`
	StuName     string `json:"stuName"`
	NickName    string `json:"nickName"`
	College     string `json:"college"`
	Major       string `json:"major"`
	Classes     string `json:"classes"`
	Birthday    string `json:"birthday"`
	Sex         string `json:"sex"`
	Nation      string `json:"nation"`
	NativePlace string `json:"nativePlace"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	QqNum       string `json:"qqNum"`
	WechatNum   string `json:"wechatNum"`
}

type CourseResp struct {
	ClassName  string `json:"className"`
	TeachClass string `json:"teachClass"`
	Teacher    string `json:"teacher"`
	StartWeek  int    `json:"startWeek"`
	EndWeek    int    `json:"endWeek"`
	Section    int    `json:"section"`
	WeekDay    int    `json:"weekDay"`
	Classroom  string `json:"classroom"`
}

type ScoreResp struct {
	SchoolTerm     string  `json:"schoolTerm"`
	CourseNum      string  `json:"courseNum"`
	CourseName     string  `json:"courseName"`
	Grade          string  `json:"grade"`
	ScoreFlag      string  `json:"scoreFlag"`
	CourseCredit   float64 `json:"courseCredit"`
	CourseHours    float64 `json:"courseHours"`
	GradePoint     float64 `json:"gradePoint"`
	EvaluationMode string  `json:"evaluationMode"`
	ExamNature     string  `json:"examNature"`
	CourseNature   string  `json:"courseNature"`
}

type GraduateStudentResp struct {
	Id         string `json:"id"`
	StudentNum string `json:"studentNum"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Degree     string `json:"degree"`
	TutorName  string `json:"tutorName"`
	Academy    string `json:"academy"`
	Specialty  string `json:"specialty"`
	//Grade      int    `json:"grade"`
	Avatar string `json:"avatar"`
}

type GraduateScoreResp struct {
	Name   string  `json:"name"`
	Credit float64 `json:"credit"`
	Term   int     `json:"term"`
	Point  string  `json:"point"`
}
