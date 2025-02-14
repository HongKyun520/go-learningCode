package web

import (
	"GoInAction/webook/internal/domain"
	"GoInAction/webook/internal/service"
	"GoInAction/webook/internal/web/jwt"
	"GoInAction/webook/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	// 创作者接口
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)

	// 获取文章详情
	g.GET("/detail/:id", h.Detail)
	// 获取文章列表（分页）
	g.POST("/list", h.List)

	pub := g.Group("/pub")
	// 读者获取文章详情
	pub.POST("/detail", h.PubDetail)

}

// 读者获取文章详情
func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "id 参数错误",
			Code: 4,
		})
		h.logger.Warn("查询文章失败，id 格式不对",
			logger.String("id", idstr),
			logger.Error(err))
		return
	}

	art, err := h.svc.GetPubById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Msg:  "系统错误",
			Code: 5,
		})
		h.logger.Error("查询文章失败，系统错误",
			logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,

			Content:    art.Content,
			AuthorId:   art.Author.Id,
			AuthorName: art.Author.Name,

			Status: art.ToUint8(),
			CTime:  art.Ctime.Format(time.DateTime),
			Utime:  art.Utime.Format(time.DateTime),
		},
	})
}

// 获取文章详情
func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	var id int64
	// 将字符串id转换为int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}
	uc := ctx.MustGet("user").(jwt.UserClaims)

	art, err := h.svc.GetById(ctx, id)
	if err != nil {
		h.logger.Error("获取文章详情失败",
			logger.Int64("article_id", id),
			logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 判断当前登录用户是否是文章作者
	if art.Author.Id != uc.Id {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.logger.Error("非法访问",
			logger.Int64("article_id", id),
			logger.Int64("user_id", uc.Id))
		return
	}

	vo := ArticleVo{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.ToUint8(),
		CTime:   art.Ctime.Format(time.DateTime),
		Utime:   art.Utime.Format(time.DateTime),
		Content: art.Content,
	}

	ctx.JSON(http.StatusOK, Result{
		Data: vo,
	})
}

// 获取文章列表（分页）
func (h *ArticleHandler) List(ctx *gin.Context) {
	var page Page

	if err := ctx.Bind(&page); err != nil {
		h.logger.Error("参数解析错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		return
	}

	// 获取当前登录用户
	uc := ctx.MustGet("user").(jwt.UserClaims)
	arts, err := h.svc.GetByAuthor(ctx, uc.Id, page.Offset, page.Limit)
	if err != nil {
		h.logger.Error("获取文章列表失败",
			logger.Int64("author_id", uc.Id),
			logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 记录日志
		h.logger.Error("获取文章列表失败",
			logger.Int64("author_id", uc.Id),
			logger.Int("offset", page.Offset),
			logger.Int("limit", page.Limit),
			logger.Error(err))
		return
	}

	vo := slice.Map[domain.Article, ArticleVo](arts, func(idx int, src domain.Article) ArticleVo {
		return ArticleVo{
			Id:       src.Id,
			Title:    src.Title,
			Abstract: src.Abstract(),
			Status:   src.ToUint8(),
			CTime:    src.Ctime.Format(time.DateTime),
			Utime:    src.Utime.Format(time.DateTime),
		}
	})
	ctx.JSON(http.StatusOK, Result{
		Data: vo,
	})
}

// 撤回文章
func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.logger.Error("内部错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "内部错误",
		})
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	err := h.svc.Withdraw(ctx, req.Id, uc.Id)
	if err != nil {
		h.logger.Error("撤回文章失败", logger.Int64("author_id", uc.Id), logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

// 接收 Article 输入，返回一个Id
func (h *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorId int64  `json:"author_id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.logger.Error("内部错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "内部错误",
		})
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	artId, err := h.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Id,
		},
	})
	if err != nil {
		h.logger.Error("保存文章失败", logger.Int64("author_id", uc.Id), logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: artId,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {

	type Req struct {
		Id       int64  `json:"id"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		AuthorId int64  `json:"author_id"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		h.logger.Error("内部错误", logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "内部错误",
		})
		return
	}

	uc := ctx.MustGet("user").(jwt.UserClaims)
	artId, err := h.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uc.Id,
		},
	})
	if err != nil {
		h.logger.Error("发表文章失败", logger.Int64("author_id", uc.Id), logger.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "success",
		Data: artId,
	})

}
