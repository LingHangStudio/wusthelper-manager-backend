package http

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"wusthelper-manager-go/app/service"
	"wusthelper-manager-go/library/ecode"
)

type TermListResp struct {
	Id        int64  `json:"id"`
	Term      string `json:"term"`
	StartDate string `json:"startDate"`
}

type TermAddReq struct {
	Term      string `json:"term" binging:"required"`
	StartDate string `json:"startDate" binging:"required"`
}

type TermModifyReq struct {
	Id        int64   `json:"id" binging:"required"`
	Term      *string `json:"term"`
	StartDate *string `json:"startDate"`
}

func getTermList(c *gin.Context) {
	terms, err := srv.GetTermList()
	if err != nil {
		responseEcode(c, err)
		return
	}

	respList := make([]TermListResp, len(*terms))
	for i, term := range *terms {
		respList[i] = TermListResp{
			Id:        term.ID,
			Term:      *term.Term,
			StartDate: term.Start.Format(_defaultDateFormat),
		}
	}

	responseData(c, respList)
}

func addTerm(c *gin.Context) {
	req := new(TermAddReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	startTime, err := time.Parse(_defaultDateFormat, req.StartDate)
	if err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	term := service.TermAddParam{
		Term:  req.Term,
		Start: startTime,
	}

	err = srv.AddTerm(&term)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

func modifyTerm(c *gin.Context) {
	req := new(TermModifyReq)
	if err := c.ShouldBindJSON(req); err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	newTerm := service.TermModifyParam{
		ID:   req.Id,
		Term: req.Term,
	}

	if req.StartDate != nil {
		startTime, err := time.Parse(_defaultDateFormat, *req.StartDate)
		if err != nil {
			responseEcode(c, ecode.ParamWrong)
			return
		}

		newTerm.Start = &startTime
	}

	err := srv.ModifyTerm(&newTerm)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}

func deleteTerm(c *gin.Context) {
	reqId, has := c.GetPostForm("id")
	if !has {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	id, err := strconv.ParseInt(reqId, 10, 64)
	if err != nil {
		responseEcode(c, ecode.ParamWrong)
		return
	}

	err = srv.DeleteTerm(id)
	if err != nil {
		responseEcode(c, err)
		return
	}

	responseData(c, nil)
}
