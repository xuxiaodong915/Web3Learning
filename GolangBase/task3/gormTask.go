package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、
// age （学生年龄，整数类型）、 grade （学生年级，字符串类型
type Student struct {
	ID    uint
	Name  string
	Age   uint8
	Grade string
}

// accounts 表（包含字段 id 主键， balance 账户余额）
type Account struct {
	ID      uint
	Balance float64
}

// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
type Transaction struct {
	ID            uint
	FromAccountID uint
	ToAccountID   uint
	Amount        float64
}

// employees 表，包含字段 id 、 name 、 department 、 salary 。
type Employee struct {
	gorm.Model
	Name       string
	Department string
	Salary     float64
}

// 有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）

// User
type User struct {
	ID           uint
	Name         string
	Posts        []Post
	ArticleCount uint
}

// Post
type Post struct {
	ID       uint
	UserID   uint
	Status   string
	PostName string
	Comments []Comment
}

// Comment
type Comment struct {
	ID      uint
	PostID  uint
	Content string
}

type PostWithCommentCount struct {
	Post
	CommentCount int64
}

func initDB() *gorm.DB {
	db := ConnectDB()
	err := db.AutoMigrate(&Student{}, &Account{}, &Transaction{}, &Employee{}, &User{}, &Post{}, &Comment{})
	if err != nil {
		panic(err)
	}
	return db
}

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:33061)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	return db
}

func Save(db *gorm.DB, student *Student) {
	result := db.Create(student)
	if result.Error != nil {
		fmt.Println("Error creating student:", result.Error)
	}
	var aff = result.RowsAffected
	fmt.Println("RowsAffected:", aff)
}

func transactionDemo(db *gorm.DB, from *Account, to *Account, amount float64) {
	err := db.Transaction(func(tx *gorm.DB) error {

		fromSaveErr := tx.Debug().Model(&from).Where("balance > ?", amount).Update("balance", gorm.Expr("balance - ?", amount)).Error
		if fromSaveErr != nil {
			return fromSaveErr
		}
		toErr := tx.Debug().Model(&to).Update("balance", gorm.Expr("balance + ?", amount)).Error
		if toErr != nil {
			return toErr
		}

		transaction := Transaction{FromAccountID: from.ID, ToAccountID: to.ID, Amount: amount}
		traErr := tx.Debug().Create(&transaction).Error
		return traErr
	})
	if err != nil {
		fmt.Println("已回滚", err)
	} else {
		fmt.Println("已提交")
	}
}

func findEmployee(db *gorm.DB, dept string) {

	var emp = []Employee{}
	db.Debug().Unscoped().Where("department = ?", dept).Find(&emp)
	fmt.Println("根据部门查结果：", emp)

	var one = Employee{}
	query := db.Unscoped().Model(&Employee{}).Select("max(salary) as salary").First(&one)
	db.Debug().Unscoped().Model(&Employee{}).Where("salary = ?", query).First(&one)
	fmt.Println("薪水最高的：", one)
}

func (p *Post) AfterCreate(db *gorm.DB) (err error) {
	return db.Model(&User{}).Where("id=?", p.UserID).UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error
}

func (c *Comment) AfterDelete(db *gorm.DB) (err error) {
	var count int64
	db.Model(&Comment{}).Where("post_id =?", c.PostID).Count(&count)
	if count == 0 {
		return db.Model(&Post{}).Where("id=?", c.PostID).Update("comment_status", "无评论").Error
	}
	return nil
}

func main() {
	// 1.SQL语句练习-基本CRUD操作

	// ALTER TABLE students CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
	db := initDB()
	student := Student{
		Name:  "张三4",
		Age:   20,
		Grade: "三年级",
	}
	//Save(db, &student)
	fmt.Print(student)

	// 2.SQL语句练习-事务语句
	from := Account{Balance: 100}
	to := Account{Balance: 200}
	db.Create(&from)
	db.Create(&to)
	transactionDemo(db, &from, &to, 50)

	// 3.Sqlx入门-使用SQL扩展库进行查询
	employee := Employee{Name: "王五", Department: "技术部", Salary: 500}
	db.Create(&employee)
	employee2 := Employee{Name: "赵六", Department: "技术部", Salary: 800}
	db.Create(&employee2)
	findEmployee(db, "技术部")
	//4.Sqlx入门-实现类型安全映射

	// 5.进阶gorm-模型定义
	user := User{
		Name: "xx",
		Posts: []Post{
			{PostName: "消失的前车", Status: "", Comments: []Comment{
				{Content: "fsd牛逼"},
				{Content: "充钱帝"},
				{Content: "我遥遥领先"},
			},
			},
			{PostName: "城区大比拼", Status: "", Comments: []Comment{
				{Content: "不服不行"},
				{Content: "怎么没人说话"},
			},
			},
		},
	}
	// 6.进阶gorm-模型定义-创建这些模型对应的数据库表
	//db.Create(&user)
	fmt.Print(user)

	// 7. 进阶gorm-查询某个用户发布的所有文章及其对应的评论信息
	var user2 User
	userID := 1
	err := db.Preload("Posts.Comments").First(&user2, userID).Error
	if err != nil {
		fmt.Println("查询失败：", err)
	} else {
		fmt.Println("用户：%s\n", user2.Name)
		for _, post := range user2.Posts {
			fmt.Printf("文章：%s\n", post.PostName)
			for _, comment := range post.Comments {
				fmt.Print(" 评论：%s\n", comment.Content)
			}
		}
	}
	fmt.Print(user2)
	// 8.进阶gorm-使用Gorm查询评论数量最多的文章信息。
	var result PostWithCommentCount
	db.Debug().Model(&Post{}).
		Select("posts.*,count(comments.id) as comment_count").
		Joins("left join comments on comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count desc").
		Limit(1).
		Scan(&result)
	fmt.Printf("评论最多的文章：%s，评论数：%d\n", result.PostName, result.CommentCount)
	// 9.钩子函数1 创建时自动更新用户的文章数量
	// 10.钩子函数2 在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
}
