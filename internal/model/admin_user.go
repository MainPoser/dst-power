package model

const (
	StateOpen  = 1
	StateClose = 0
)

type AdminUser struct {
	Model
	RealName  string `gorm:"column:real_name;NOT NULL" json:"realName"`              // 真实姓名
	Username  string `gorm:"column:username;NOT NULL" json:"username"`               // 登录用户名
	Gender    int    `gorm:"column:gender;default:0;NOT NULL" json:"gender"`         // 性别:1男 2女 3保密
	Avatar    string `gorm:"column:avatar;NOT NULL" json:"avatar"`                   // 头像
	Mobile    string `gorm:"column:mobile;NOT NULL" json:"mobile"`                   // 手机号
	Email     string `gorm:"column:email;NOT NULL" json:"email"`                     // 邮箱地址
	Address   string `gorm:"column:address;NOT NULL" json:"address"`                 // 地址
	Password  string `gorm:"column:password;NOT NULL" json:"password"`               // 登录密码
	Salt      string `gorm:"column:salt;NOT NULL" json:"salt"`                       // 盐加密
	Intro     string `gorm:"column:intro" json:"intro"`                              // 备注
	Status    int    `gorm:"column:status;default:1;NOT NULL" json:"status"`         // 状态：1正常 2禁用
	LoginNum  int    `gorm:"column:login_num;default:0;NOT NULL" json:"login_num"`   // 登录次数
	LoginIp   string `gorm:"column:login_ip;NOT NULL" json:"login_ip"`               // 最近登录ip
	LoginTime int64  `gorm:"column:login_time;default:0;NOT NULL" json:"login_time"` // 最近登录时间
}

func (a *AdminUser) TableName() string {
	return "admin_user"
}
