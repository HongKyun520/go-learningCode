package integration

import (
	"GoInAction/webook/internal/integration/startup"
	"GoInAction/webook/internal/repository/dao"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ArticleHandlerSuite struct {
	suite.Suite
	db     *gorm.DB
	server *gin.Engine
}

// 在每个测试用例结束后，清空文章表
func (s *ArticleHandlerSuite) TearDownTest() {
	s.db.Exec("truncate TABLE articles")
}

func (s *ArticleHandlerSuite) SetupSuite() {
	// TODO: 解决server的初始化问题
	s.server = gin.Default()
	s.db = startup.InitDB()
}

func TestArticleHandler(t *testing.T) {
	suite.Run(t, &ArticleHandlerSuite{})
}

func (s *ArticleHandlerSuite) TestArticleHandler_Edit() {

	t := s.T()
	testCases := []struct {
		name       string
		before     func(*testing.T)
		after      func(*testing.T)
		req        Article
		wantCode   int
		wantResult Result[int64]
	}{
		{
			name:   "新建帖子",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				// 验证数据库已经插入了一条数据
				var article dao.Article
				s.db.Where("author_id = ?", 1).First(&article)
				assert.True(t, article.Ctime > 0)
				assert.True(t, article.Utime > 0)
				assert.True(t, article.Id > 0)
				assert.Equal(t, "我的标题", article.Title)
				assert.Equal(t, "我的内容", article.Content)
				assert.Equal(t, int64(1), article.AuthorId)
			},
			req: Article{
				Title:    "我的标题",
				Content:  "我的内容",
				AuthorId: 1,
			},
			wantCode: http.StatusOK,
			wantResult: Result[int64]{
				Code: 0,
				Msg:  "success",
				Data: 1,
			},
		},
		{
			name: "修改帖子",
			before: func(t *testing.T) {
				err := s.db.Create(dao.Article{
					Id:       2,
					Title:    "我的标题",
					Content:  "我的内容",
					AuthorId: 123,
					Status:   1,
					Ctime:    456,
					Utime:    789,
				}).Error
				assert.NoError(t, err)
			},
			req: Article{
				Id:       2,
				Title:    "新的标题",
				Content:  "新的内容",
				AuthorId: 123,
			},
			after: func(t *testing.T) {
				var article dao.Article
				err := s.db.Where("id = ?", 2).First(&article).Error
				assert.NoError(t, err)
				assert.True(t, article.Utime > 789)
				article.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    456,
					Status:   1,
				}, article)
			},
			wantCode: http.StatusOK,
			wantResult: Result[int64]{
				Code: 0,
				Msg:  "success",
				Data: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)

			reqBody, err := json.Marshal(tc.req)
			assert.NoError(t, err)

			// 准备req和记录
			req, err := http.NewRequest(http.MethodPost,
				"/articles/edit",
				bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			s.server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var res Result[int64]
			err = json.NewDecoder(recorder.Body).Decode(&res)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, res)

		})
	}
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type Article struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorId int64  `json:"author_id"`
}
